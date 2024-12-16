package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// URL struct to store original and shortened URLs
type URL struct {
	ID           string `json:"id"`
	OriginalURL  string `json:"original_url"`
	ShortURL     string `json:"short_url"`
	CreationDate string `json:"creation_date"`
}

// In-memory database to store URLs
var urldb = make(map[string]URL)

// Generate a short URL using MD5 hash (first 8 characters)
func generateShortURL(originalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:8]
}

// Create a new short URL and store it in the database
func createURL(originalURL string) string {
	shortID := generateShortURL(originalURL)
	shortURL := "http://localhost:3000/redirect/" + shortID
	url := URL{
		ID:           shortID,
		OriginalURL:  originalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now().Format(time.RFC3339),
	}
	urldb[shortID] = url
	return shortURL
}

// Retrieve the original URL using the short ID
func getURL(shortID string) (URL, error) {
	url, exists := urldb[shortID]
	if !exists {
		return URL{}, fmt.Errorf("URL not found with ID: %s", shortID)
	}
	return url, nil
}

// Middleware to handle CORS for frontend integration
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func isFrontendRequest(r *http.Request) bool {
    origin := r.Header.Get("Origin")
    if origin == "https://friendly-memory-44vw69q44wq2j7j-3001.app.github.dev" {
        return true
    }
    return false
}

// Handler for shortening URLs
func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Shorten url endpoint")
	fmt.Println("Shortening url function called")

	if !isFrontendRequest(r) {
        http.Error(w, "Forbidden: Request not allowed", http.StatusForbidden)
        return
    }

	if r.Method == http.MethodOptions {
		enableCORS(w)
		w.WriteHeader(http.StatusOK) // Respond with a 200 OK
		return
	}
	enableCORS(w)
	

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.URL == "" {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Ensure the URL has a valid format
	originalURL := requestData.URL
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		http.Error(w, "URL must start with http:// or https://", http.StatusBadRequest)
		return
	}

	shortURL := createURL(originalURL)

	responseData := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

// Handler for redirecting short URLs to original URLs
func redirectURLHandler(w http.ResponseWriter, r *http.Request) {

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		enableCORS(w)
		w.WriteHeader(http.StatusOK) // Respond with a 200 OK
		return
	}

	// Enable CORS for this handler
	enableCORS(w)

	// Extract short URL ID
	shortID := strings.TrimPrefix(r.URL.Path, "/redirect/")
	url, err := getURL(shortID)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

// Root handler (for testing purposes)
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		enableCORS(w)
		w.WriteHeader(http.StatusOK) // Respond with a 200 OK
		return
	}

	// Enable CORS for this handler
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	// Regular processing for other methods (like GET)
	fmt.Println("Root handler called")
}

// Main function to start the server
func main() {
	fmt.Println("Starting URL shortener server...")

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/shorten", shortenURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	port := ":3000"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

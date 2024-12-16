package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// url shortner detailed struct
type URL struct {
	ID           string `json:"id"`
	OriginalURL  string `json:"original_url"`
	ShortURL     string `json:"short_url"`
	CreationDate string `json:"creation_date"`
}

// in memory database
/*
 123 --> {
     ID:           "123",
     OriginalURL:  "https://www.example.com",
     ShortURL:     "123",
     CreationDate: "time.now()",
  }

*/
var urldb = make(map[string]URL)

func generateShortURL(OriginalUrl string) string {

	hasher := md5.New()
	hasher.Write([]byte(OriginalUrl))
	// fmt.Println("hasher: ", hasher)
	data := hasher.Sum(nil)
	// fmt.Println("hasher data: ", data)
	hash := hex.EncodeToString(data)
	fmt.Println("hasher data: ", data)
	// fmt.Println("hash: ", hash)
	fmt.Println("final hash: ", hash[:8])
	return hash[:8]

}

func createURL(originalURL string) string {
	shortURL := generateShortURL(originalURL)
	url := URL{
		ID:           shortURL,
		OriginalURL:  originalURL,
		ShortURL:     "http://localhost:8080/" + shortURL,
		CreationDate: "time.now()",
	}
	urldb[url.ID] = url
	fmt.Println("URL created: ", url)
	// fmt.Println("URL DB: ", urldb)
	return url.ShortURL
}

func getURL(id string) (URL, error) {
	url, ok := urldb[id]
	if !ok {
		return URL{}, fmt.Errorf("URL not found with ID: %s", id)
	}
	return url, nil
}

func redirectURL(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortURL := createURL(data.URL)
	// fmt.Fprintf(w,shortURL)

	response := struct {
		ShortedURL string `json:"shortedURL"`
	}{ShortedURL: shortURL}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get method")
	fmt.Fprintf(w, "Hello from server")

}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusSeeOther)
}

func main() {
	fmt.Println("Url shortener...")
	generateShortURL("http://example.com/221212")

	// Register the handler function to handle all requests to the root URL ("/")

	http.HandleFunc("/", handler)
	http.HandleFunc("/shorten", redirectURL)
	http.HandleFunc("/redirect/", redirectURLHandler)

	// Start the server on the default port 3000
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("Server started on port 3000")
}

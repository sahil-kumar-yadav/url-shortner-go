package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
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
// var urldb = make(map[string]URL)

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
func main() {
	fmt.Println("Url shortener...")
	generateShortURL("http://example.com/221212")
}

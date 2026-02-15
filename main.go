package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title" json:"title"`
	Link        string `xml:"link" json:"link"`
	PubDate     string `xml:"pubDate" json:"pub_date"`
	Description string `xml:"description" json:"description"`
}

func generateGoogleQuery(query, source string) string {
	return fmt.Sprintf("https://news.google.com/rss/search?q=%s+site:%s", query, source)
}

func getSearchResult(query string) ([]Item, error) {
	sources := [4]string{
		"bbc.co.uk",
		"g1.globo.com",
		"aljazeera.com",
		"reuters.com",
	}

	results := []Item{}

	for _, source := range sources {
		resp, err := http.Get(generateGoogleQuery(query, source))
		if err != nil {
			fmt.Printf("Error fetching news: %v\n", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			fmt.Printf("Error reading body: %v\n", err)
			return results, errors.New("Error reading body")
		}

		var rss RSS
		if err := xml.Unmarshal(body, &rss); err != nil {
			fmt.Printf("Error parsing XML: %v\n", err)
			return results, errors.New("Error parsing XML")
		}

		if len(rss.Channel.Items) > 0 {
			results = append(results, rss.Channel.Items[0])
		}
	}

	return results, nil
}

func handleGetSearchResult(w http.ResponseWriter, r *http.Request) {
	query := r.PathValue("query")

	if query == "" {
		http.Error(w, "Missing query", http.StatusBadRequest)
		return
	}

	res, err := getSearchResult(query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Returned data: %v\n", res)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{query}", handleGetSearchResult)
	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

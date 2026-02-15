package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
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
			return results, errors.New("Error fetching news")
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
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
		} else {
			return results, errors.New("No items found.")
		}
	}

	return results, nil
}

func main() {

	if len(os.Args) < 3 {
		fmt.Printf("usage: go run main.go <query>")
	}

	res, err := getSearchResult(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %e", err)
	}

	fmt.Printf("Returned data: %v", res)

}

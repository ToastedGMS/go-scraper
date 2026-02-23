package main

import (
	"io"
	"log"
	"net/http"
)

func metro() {
	req, _ := http.NewRequest("GET", "https://www.metropoles.com/wp-json/wp/v2/posts?search=kanye&_embed&per_page=1", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "https://www.metropoles.com/")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalf("An error ocurred %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	log.Printf(sb)

}

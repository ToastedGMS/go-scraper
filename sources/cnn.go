package main

import (
	"io"
	"log"
	"net/http"
)

func cnn() {
	req, err := http.NewRequest(http.MethodGet, "https://admin.cnnbrasil.com.br/wp-json/wp/v2/posts?search=Kanye&_embed&per_page=1", nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	req.Header.Add("Accept", "application/json")

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

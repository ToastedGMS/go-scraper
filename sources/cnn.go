package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func cnn(query string) {
	params := url.Values{}
	params.Add("search", query)
	params.Add("_embed", "")
	params.Add("per_page", "1")

	baseUrl := "https://admin.cnnbrasil.com.br/wp-json/wp/v2/posts"
	fullUrl := baseUrl + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
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

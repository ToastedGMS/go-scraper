package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func g1() {
	payload := []map[string]interface{}{
		{
			"search_profile": "sp_g1_globo_com",
			"query":          "g1.info_query_recency",
			"params": map[string]interface{}{
				"q":    "kanye",
				"from": 0,
				"size": 1,
			},
		}}

	postBody, _ := json.Marshal(payload)

	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest(http.MethodPost, "https://busca.globo.com/v1/search", responseBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-tenant-id", "g1")
	req.Header.Add("Origin", "https://g1.globo.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)

}

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func main() {
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

	resp, err := http.NewRequest(http.MethodPost, "https://busca.globo.com/v1/search", responseBody)
	resp.Header.Add("Content-Type", "application/json")
	resp.Header.Add("x-tenant-id", "g1")
	resp.Header.Add("Origin", "https://g1.globo.com")

	response, err := http.DefaultClient.Do(resp)

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

package sources

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func G1(query string) {
	type ParsedG1Response []struct {
		Result struct {
			Hits struct {
				Hits []struct {
					Source struct {
						Title     string `json:"title"`
						Publisher string `json:"publisher"`
						Issued    string `json:"issued"`
						URL       string `json:"url"`
						Thumbnail string `json:"thumbnail"`
					} `json:"_source"`
				} `json:"hits"`
			} `json:"hits"`
		} `json:"result"`
	}

	payload := []map[string]interface{}{
		{
			"search_profile": "sp_g1_globo_com",
			"query":          "g1.info_query_recency",
			"params": map[string]interface{}{
				"q":    query,
				"from": 0,
				"size": 1,
			},
		}}

	postBody, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error formatting request payload for G1: %v", err)
	}

	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest(http.MethodPost, "https://busca.globo.com/v1/search", responseBody)
	if err != nil {
		log.Fatalf("Error creating request to G1: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-tenant-id", "g1")
	req.Header.Add("Origin", "https://g1.globo.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer response.Body.Close()

	var parsed ParsedG1Response

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response body from G1: %v", err)
	}

	if err := json.Unmarshal(body, &parsed); err != nil {
		log.Fatalf("Error parsing response body from G1: %v", err)
	}

	log.Printf("%+v", parsed[0].Result.Hits.Hits[0].Source)

}

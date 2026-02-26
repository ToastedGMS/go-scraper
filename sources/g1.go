package sources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ToastedGMS/go-scraper/types"
)

func G1(query string) (types.Article, error) {
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
		return types.Article{}, fmt.Errorf("G1: payload marshalling error: %w", err)
	}

	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest(http.MethodPost, "https://busca.globo.com/v1/search", responseBody)
	if err != nil {
		return types.Article{}, fmt.Errorf("Error creating request to G1: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-tenant-id", "g1")
	req.Header.Add("Origin", "https://g1.globo.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return types.Article{}, fmt.Errorf("Error sending request to G1: %w", err)
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return types.Article{}, fmt.Errorf("Unexpected response status from G1: %d", response.StatusCode)
	}

	var parsed ParsedG1Response

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return types.Article{}, fmt.Errorf("Error reading response from G1: %w", err)
	}

	if err := json.Unmarshal(body, &parsed); err != nil {
		return types.Article{}, fmt.Errorf("Error parsing response from G1: %w", err)
	}

	var final types.Article

	if len(parsed) == 0 || len(parsed[0].Result.Hits.Hits) == 0 {
		return types.Article{}, fmt.Errorf("No results found for query: %s", query)
	}

	hit := parsed[0].Result.Hits.Hits[0].Source

	final.Title = hit.Title
	final.Date = hit.Issued
	final.Img = hit.Thumbnail
	final.Source = hit.Publisher
	final.URL = hit.URL

	return final, nil

}

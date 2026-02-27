package sources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ToastedGMS/go-scraper/types"
)

func G1(query string) ([]types.Article, error) {
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
				"size": 5,
			},
		}}

	postBody, err := json.Marshal(payload)
	if err != nil {
		return []types.Article{}, fmt.Errorf("G1: payload marshalling error: %w", err)
	}

	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest(http.MethodPost, "https://busca.globo.com/v1/search", responseBody)
	if err != nil {
		return []types.Article{}, fmt.Errorf("Error creating request to G1: %w", err)
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-tenant-id", "g1")
	req.Header.Set("Origin", "https://g1.globo.com")
	req.Header.Set("Referer", "https://g1.globo.com/")

	response, err := Client.Do(req)
	if err != nil {
		return []types.Article{}, fmt.Errorf("G1 request failed or timed out: %w", err)
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return []types.Article{}, fmt.Errorf("Unexpected response status from G1: %d", response.StatusCode)
	}

	var parsed ParsedG1Response

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []types.Article{}, fmt.Errorf("Error reading response from G1: %w", err)
	}

	if err := json.Unmarshal(body, &parsed); err != nil {
		return []types.Article{}, fmt.Errorf("Error parsing response from G1: %w", err)
	}

	var final []types.Article

	if len(parsed) == 0 || len(parsed[0].Result.Hits.Hits) == 0 {
		return []types.Article{}, fmt.Errorf("No results found for query: %s", query)
	}
	for _, item := range parsed[0].Result.Hits.Hits {
		hit := item.Source

		if hit.Title != "" && hit.URL != "" {
			final = append(final, types.Article{
				Title:  hit.Title,
				URL:    hit.URL,
				Img:    hit.Thumbnail,
				Date:   hit.Issued,
				Source: hit.Publisher,
			})
		}
	}

	return final, nil
}

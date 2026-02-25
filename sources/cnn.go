package sources

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/ToastedGMS/go-scraper/types"
)

func Cnn(query string) types.Article {
	type ParsedCnnResponse []struct {
		Date             string `json:"date"`
		Link             string `json:"link"`
		FeaturedMediaURL string `json:"jetpack_featured_media_url"`
		Title            struct {
			Rendered string `json:"rendered"`
		} `json:"title"`
		Source string
	}
	params := url.Values{}
	params.Add("search", query)
	params.Add("per_page", "1")

	baseUrl := "https://admin.cnnbrasil.com.br/wp-json/wp/v2/posts"
	fullUrl := baseUrl + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		log.Fatalf("An error ocurred %v", err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An error ocurred %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("An error ocurred %v", err)
	}

	var parsed ParsedCnnResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		log.Fatalf("An error ocurred %v", err)

	}
	for i := range parsed {
		parsed[i].Source = "Cnn Brasil"
	}

	var final types.Article

	if len(parsed) == 0 {
		log.Fatalf("An error ocurred %v", err)
		return final
	}

	final.Title = parsed[0].Title.Rendered
	final.Date = parsed[0].Date
	final.Img = parsed[0].FeaturedMediaURL
	final.Source = parsed[0].Source
	final.URL = parsed[0].Link

	return final

}

package sources

import (
	"encoding/json"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/ToastedGMS/go-scraper/types"
)

func getFirstImage(htmlContent string) string {
	re := regexp.MustCompile(`src="([^"]+\.(?:jpg|jpeg|png|webp|gif))[^"]*"`)
	match := re.FindStringSubmatch(htmlContent)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func Metro(query string) types.Article {
	type ParsedMetroResponse []struct {
		Date    string `json:"date"`
		Link    string `json:"link"`
		Img     string
		Content struct {
			Rendered string `json:"rendered"`
		} `json:"content"`
		Title struct {
			Rendered string `json:"rendered"`
		} `json:"title"`
		Source string
	}

	params := url.Values{}
	params.Add("search", query)
	params.Add("per_page", "1")

	baseUrl := "https://www.metropoles.com/wp-json/wp/v2/posts"
	fullUrl := baseUrl + "?" + params.Encode()

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		log.Fatalf("An error ocurred %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "https://www.metropoles.com/")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An error ocurred %v", err)
	}

	defer resp.Body.Close()

	var parsed ParsedMetroResponse

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("An error ocurred %v", err)
	}

	if err := json.Unmarshal(body, &parsed); err != nil {
		log.Fatalf("An error ocurred %v", err)

	}
	for i := range parsed {
		parsed[i].Source = "Metropoles"
		parsed[i].Img = getFirstImage(parsed[i].Content.Rendered)
		parsed[i].Title.Rendered = html.UnescapeString(parsed[i].Title.Rendered)
		parsed[i].Content.Rendered = ""
	}

	var final types.Article

	if len(parsed) == 0 {
		log.Fatalf("An error ocurred %v", err)
		return final
	}

	final.Title = parsed[0].Title.Rendered
	final.Img = parsed[0].Img
	final.URL = parsed[0].Link
	final.Date = parsed[0].Date
	final.Source = parsed[0].Source

	return final

}

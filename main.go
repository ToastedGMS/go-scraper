package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/ToastedGMS/go-scraper/sources"
	"github.com/ToastedGMS/go-scraper/types"
)

func RunScrapers(query string) ([]types.Article, []error) {
	var stack sync.WaitGroup

	type Results struct {
		article []types.Article
		errors  error
	}

	responseChannel := make(chan Results, 3)

	functions := []func(string) ([]types.Article, error){
		sources.Cnn,
		sources.G1,
		sources.Metro,
	}
	stack.Add(len(functions))

	for _, function := range functions {
		go func(f func(string) ([]types.Article, error)) {
			defer stack.Done()
			article, errors := f(query)
			responseChannel <- Results{article, errors}
		}(function)
	}

	go func() {
		stack.Wait()
		close(responseChannel)
	}()

	var results []types.Article
	var issues []error

	for result := range responseChannel {
		if result.errors != nil {
			issues = append(issues, result.errors)
			continue
		}
		results = append(results, result.article...)
	}

	if len(issues) > 0 {
		for _, err := range issues {
			log.Printf("Scraper Warning: %v", err)
		}
	}

	return results, issues
}

func ScraperHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	if query == "" {
		http.Error(w, "No query provided", http.StatusBadRequest)
		return
	}

	articles, errors := RunScrapers(query)

	w.Header().Set("Content-Type", "application/json")
	if len(articles) == 0 && len(errors) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Failed to fetch data"})
		return
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Date > articles[j].Date
	})

	var stories = GroupArticles(articles)

	sort.Slice(stories, func(i, j int) bool {
		return len(stories[i].Articles) > len(stories[j].Articles)
	})
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	encoder.Encode(stories)
}

func processString(val string) []string {
	val = strings.ToLower(val)
	f := func(r rune) rune {
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			return -1
		}
		return r
	}
	val = strings.Map(f, val)

	stringSlice := strings.Split(val, " ")

	var filteredStringSlice []string
	for _, item := range stringSlice {
		item = strings.TrimSpace(item)
		if utf8.RuneCountInString(item) >= 4 {
			filteredStringSlice = append(filteredStringSlice, item)
		}
	}
	return filteredStringSlice
}

func countEqualStrings(slice1, slice2 []string) int {
	count := 0

	for _, word := range slice1 {
		for _, word2 := range slice2 {
			if word == word2 {
				count++
				break
			}
		}
	}
	return count
}

func GroupArticles(articles []types.Article) []types.Story {
	ArticleMap := make(map[int]bool)
	var stories []types.Story

	for i := 0; i < len(articles); i++ {
		ArticleMap[i] = false
	}

	for i := 0; i < len(articles); i++ {
		if ArticleMap[i] == true {
			continue
		}
		var story types.Story

		var longestTitle string
		story.Headline = articles[i].Title
		story.Articles = append(story.Articles, articles[i])
		ArticleMap[i] = true

		for j := 0; j < len(articles); j++ {
			if ArticleMap[j] == true {
				continue
			}

			for _, article := range story.Articles {
				if countEqualStrings(processString(article.Title), processString(articles[j].Title)) >= 3 {
					story.Articles = append(story.Articles, articles[j])
					ArticleMap[j] = true
					break
				}
			}
		}

		for _, article := range story.Articles {
			if utf8.RuneCountInString(article.Title) > utf8.RuneCountInString(longestTitle) {
				longestTitle = article.Title
			}
			story.Headline = longestTitle
		}
		stories = append(stories, story)

	}
	return stories
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/search", ScraperHandler)

	port := "8080"
	log.Printf("Server running on port: %v", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

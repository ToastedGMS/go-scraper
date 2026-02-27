package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

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
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	encoder.Encode(articles)
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

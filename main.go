package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/ToastedGMS/go-scraper/sources"
	"github.com/ToastedGMS/go-scraper/types"
)

func RunScrapers(query string) ([]types.Article, []error) {
	var stack sync.WaitGroup
	stack.Add(3)

	type Results struct {
		article types.Article
		errors  error
	}

	responseChannel := make(chan Results, 3)

	functions := []func(string) (types.Article, error){
		sources.Cnn,
		sources.G1,
		sources.Metro,
	}

	for _, function := range functions {
		go func(f func(string) (types.Article, error)) {
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
		results = append(results, result.article)
	}

	return results, issues
}

func main() {
	fmt.Println(RunScrapers(os.Args[1]))
}

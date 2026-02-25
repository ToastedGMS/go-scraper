package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/ToastedGMS/go-scraper/sources"
	"github.com/ToastedGMS/go-scraper/types"
)

func RunScrapers(query string) []types.Article {
	var stack sync.WaitGroup
	stack.Add(3)
	channel := make(chan types.Article, 3)

	go func() {
		defer stack.Done()
		channel <- sources.G1(query)
	}()

	go func() {
		defer stack.Done()
		channel <- sources.Cnn(query)
	}()

	go func() {
		defer stack.Done()
		channel <- sources.Metro(query)
	}()

	go func() {
		stack.Wait()
		close(channel)
	}()

	var results []types.Article

	for item := range channel {
		results = append(results, item)
	}
	return results
}

func main() {
	fmt.Println(RunScrapers(os.Args[1]))
}

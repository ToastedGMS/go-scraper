# Go Scraper

Go Scraper (awaiting a better name) is a project built with the goal of further improving my knowledge and experience with Golang, by putting in place some concepts that I skipped in my previous [Contact Book](https://github.com/ToastedGMS/go-contact-book) project.

The idea is to access external APIs of popular news websites/outlets, fetch news articles related to a provided search query, compare them between each other to group related news stories, and then present this organized list for the [frontend](https://github.com/ToastedGMS/go-scraper-fe) to display.

## Why Golang?

Between Node.Js and Golang, through research, I realized each language has it's strengths and weaknesses. As Node has a better time handling I/O bound operations (such as fetching data from an API), while Go has a better time handling CPU bound operations (such as the article sorting algorithm). Besides, due to Go's concurrency model, it actually ends up performing better for handling the multiple API operations at the same time than Node would.

## Finding the APIs

To find the APIs I used, I had to analyze the network requests from each outlet's website, and look through the responses to find API responses. Through that I was able to identify the sources of articles, and then replicating the headers my browser sent them, I was able to get the direct JSON response these APIs serve. After that all I needed was to pass these responses through a normalizing layer, to have all the articles from all sources under the same format.

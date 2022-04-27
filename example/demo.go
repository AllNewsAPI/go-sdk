package main

import (
	"fmt"
	"log"
	"time"

	"github.com/FreeNews-API/go-sdk"
)

func main() {
	// Create a new client with your API key
	client, err := freenewsapi.NewClient("your-api-key")
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Example 1: Simple search
	fmt.Println("EXAMPLE 1: Simple search for 'bitcoin'")
	response, err := client.Search(&freenewsapi.SearchOptions{
		Query: "bitcoin",
		Max:   3,
	})
	if err != nil {
		log.Fatalf("Error searching: %v", err)
	}

	fmt.Printf("Found %d articles\n", response.TotalArticles)
	for _, article := range response.Articles {
		fmt.Printf("Title: %s\n", article.Title)
		fmt.Printf("Source: %s\n", article.Source.Name)
		fmt.Printf("URL: %s\n", article.URL)
		fmt.Println("---")
	}
	fmt.Println()

}
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/AllNewsAPI/go-sdk"
)

func main() {
	// Create a new client with your API key
	client, err := allnewsapi.NewClient("your-api-key")
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Example 1: Simple search
	fmt.Println("EXAMPLE 1: Simple search for 'bitcoin'")
	response, err := client.Search(&allnewsapi.SearchOptions{
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

	// Example 2: Get headlines by category
	fmt.Println("EXAMPLE 2: Get technology headlines")
	headlines, err := client.Headlines(&allnewsapi.SearchOptions{
		Category: []string{"technology"},
		Max:      3,
	})
	if err != nil {
		log.Fatalf("Error getting headlines: %v", err)
	}

	fmt.Printf("Found %d headlines\n", headlines.TotalArticles)
	for _, article := range headlines.Articles {
		fmt.Printf("Title: %s\n", article.Title)
		fmt.Printf("Source: %s\n", article.Source.Name)
		fmt.Printf("URL: %s\n", article.URL)
		fmt.Println("---")
	}
}
# FreeNewsAPI Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/FreeNews-API/go-sdk.svg)](https://pkg.go.dev/github.com/FreeNews-API/go-sdk)
[![License](https://img.shields.io/github/license/FreeNews-API/go-sdk.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/FreeNews-API/go-sdk)

The official **FreeNewsAPI SDK** for the Golang programming language. 

Fetch real-time and historical news articles and headlines from multiple sources around the world.


---

## Installation

```bash
go get github.com/FreeNews-API/go-sdk
```

---

## Usage

### Basic Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/FreeNews-API/go-sdk"
)

func main() {
	// Create a new client
	client, err := freenewsapi.NewClient("your-api-key")
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Simple search
	results, err := client.Search(&freenewsapi.SearchOptions{
		Query: "bitcoin",
		Max:   10,
	})
	if err != nil {
		log.Fatalf("Error searching: %v", err)
	}

	// Print results
	fmt.Printf("Found %d articles\n", results.TotalArticles)
	for _, article := range results.Articles {
		fmt.Printf("Title: %s\n", article.Title)
		fmt.Printf("Source: %s\n", article.Source.Name)
		fmt.Println("---")
	}
}
```

---

### Advanced Example

```go
// Advanced search with multiple parameters
includeContent := true
results, err := client.Search(&freenewsapi.SearchOptions{
	Query:      "AI startups",
	Lang:       []string{"en", "fr"},
	Category:   []string{"technology"},
	Max:        10,
	SortBy:     "relevance",
	Content:    &includeContent,
	Attributes: []string{"title", "description"},
})
if err != nil {
	log.Fatalf("Error searching: %v", err)
}
```

---

## API Reference

### Client Methods

#### `NewClient(apiKey string, options ...ClientOption) (*Client, error)`

Creates a new Free News API client.

- `apiKey`: Your Free News API key
- `options`: Optional client configuration (e.g., custom HTTP client)

---

#### `Search(options *SearchOptions) (*SearchResponse, error)`

Search for news articles with various options. 🔗 [See API Documentation](https://freenewsapi.com/documentation#search-endpoint)  

---

### SearchOptions

| Parameter    | Type                  | Description |
|--------------|------------------------|-------------|
| `q`          | `string`                | Keywords to search for |
| `startDate`  | `string`                | Start date (`YYYY-MM-DD` or `YYYY-MM-DD HH:MM:SS`) |
| `endDate`    | `string`                | End date (`YYYY-MM-DD` or `YYYY-MM-DD HH:MM:SS`) |
| `content`    | `*bool`                 | Whether to include full content |
| `lang`       | `string` or `[]string`   | Language(s) to filter by |
| `country`    | `string` or `[]string`   | Country/countries to filter by (ISO 3166 codes) |
| `region`     | `string` or `[]string`   | Region(s) to filter by |
| `category`   | `string` or `[]string`   | Category/categories to filter by |
| `max`        | `int`                   | Maximum number of results (1–100) |
| `attributes` | `string` or `[]string`   | Attributes to search in (`title`, `description`, `content`) |
| `page`       | `int`                   | Page number for pagination |
| `sortby`     | `string`                | Sort by `'publishedAt'` or `'relevance'` |
| `publisher`  | `string` or `[]string`   | Filter by publisher(s) |
| `format`     | `string`                | Response format (`json`, `csv`, `xlsx`) |

---

## Error Handling

All client methods return an `error` as the second return value.  
Always check for errors before accessing the response data.

---

## License

The MIT License (MIT). Please see the [License File](LICENSE) for more information.
// Package allnewsapi provides a client for the AllNewsAPI.
package allnewsapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is a AllNewsAPI client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Article represents a news article returned by the API.
type Article struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Content     string    `json:"content"`
	Country     string    `json:"country"`
	Region      string    `json:"region"`
	Lang        string    `json:"lang"`
	Sentiment   string    `json:"sentiment"`
	URL         string    `json:"url"`
	Image       string    `json:"image"`
	PublishedAt time.Time `json:"publishedAt"`
	Source      struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"source"`
}

// SearchResponse represents the response from the search endpoint.
type SearchResponse struct {
	TotalArticles int       `json:"totalArticles"`
	CurrentPage   int       `json:"currentPage"`
	NextPage      *int       `json:"nextPage"`
	Articles      []Article `json:"articles"`
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithTimeout sets a custom timeout for HTTP requests.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// NewClient creates a new AllNewsAPI client.
func NewClient(apiKey string, options ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("API key is required")
	}

	client := &Client{
		apiKey:  apiKey,
		baseURL: "https://api.allnewsapi.com",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Apply options
	for _, option := range options {
		option(client)
	}

	return client, nil
}

// SearchOptions contains all possible parameters for the search endpoint.
type SearchOptions struct {
	Query       string      // Search query
	StartDate   interface{} // string or time.Time
	EndDate     interface{} // string or time.Time
	Content     *bool       // Whether to include full content
	Lang        []string    // Languages to filter by
	Country     []string    // Countries to filter by
	Region      []string    // Regions to filter by
	Category    []string    // Categories to filter by
	Max         int         // Maximum number of results (1-100)
	Attributes  []string    // Attributes to search in (title, description, content)
	Page        int         // Page number for pagination
	SortBy      string      // Sort by 'publishedAt' or 'relevance'
	Publisher   []string    // Publishers to filter by
	Format      string      // Response format (json, csv, xlsx)
}

// Search searches for news articles.
func (c *Client) Search(options *SearchOptions) (*SearchResponse, error) {
	params := url.Values{}

	// Add API key
	params.Add("apikey", c.apiKey)

	// Add query parameters if provided
	if options != nil {
		if options.Query != "" {
			params.Add("q", options.Query)
		}

		// Handle start date
		if options.StartDate != nil {
			var startDate string
			switch v := options.StartDate.(type) {
			case string:
				startDate = v
			case time.Time:
				startDate = v.Format(time.RFC3339)
			default:
				return nil, errors.New("startDate must be string or time.Time")
			}
			params.Add("startDate", startDate)
		}

		// Handle end date
		if options.EndDate != nil {
			var endDate string
			switch v := options.EndDate.(type) {
			case string:
				endDate = v
			case time.Time:
				endDate = v.Format(time.RFC3339)
			default:
				return nil, errors.New("endDate must be string or time.Time")
			}
			params.Add("endDate", endDate)
		}

		// Handle boolean content parameter
		if options.Content != nil {
			if *options.Content {
				params.Add("content", "true")
			} else {
				params.Add("content", "false")
			}
		}

		// Handle array parameters
		if len(options.Lang) > 0 {
			params.Add("lang", strings.Join(options.Lang, ","))
		}
		if len(options.Country) > 0 {
			params.Add("country", strings.Join(options.Country, ","))
		}
		if len(options.Region) > 0 {
			params.Add("region", strings.Join(options.Region, ","))
		}
		if len(options.Category) > 0 {
			params.Add("category", strings.Join(options.Category, ","))
		}
		if len(options.Attributes) > 0 {
			params.Add("attributes", strings.Join(options.Attributes, ","))
		}
		if len(options.Publisher) > 0 {
			params.Add("publisher", strings.Join(options.Publisher, ","))
		}

		// Handle integer parameters
		if options.Max > 0 {
			params.Add("max", fmt.Sprintf("%d", options.Max))
		}
		if options.Page > 0 {
			params.Add("page", fmt.Sprintf("%d", options.Page))
		}

		// Handle other string parameters
		if options.SortBy != "" {
			params.Add("sortby", options.SortBy)
		}
		if options.Format != "" {
			params.Add("format", options.Format)
		}
	}

	// Build request URL
	searchURL := fmt.Sprintf("%s/v1/search?%s", c.baseURL, params.Encode())

	// Make the request
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Check for error responses
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, body)
	}

	// Parse the response
	var searchResponse SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &searchResponse, nil
}

// Headlines fetches news headlines.
func (c *Client) Headlines(options *SearchOptions) (*SearchResponse, error) {
	params := url.Values{}

	// Add API key
	params.Add("apikey", c.apiKey)

	// Add query parameters if provided
	if options != nil {
		if options.Query != "" {
			params.Add("q", options.Query)
		}

		// Handle start date
		if options.StartDate != nil {
			var startDate string
			switch v := options.StartDate.(type) {
			case string:
				startDate = v
			case time.Time:
				startDate = v.Format(time.RFC3339)
			default:
				return nil, errors.New("startDate must be string or time.Time")
			}
			params.Add("startDate", startDate)
		}

		// Handle end date
		if options.EndDate != nil {
			var endDate string
			switch v := options.EndDate.(type) {
			case string:
				endDate = v
			case time.Time:
				endDate = v.Format(time.RFC3339)
			default:
				return nil, errors.New("endDate must be string or time.Time")
			}
			params.Add("endDate", endDate)
		}

		// Handle boolean content parameter
		if options.Content != nil {
			if *options.Content {
				params.Add("content", "true")
			} else {
				params.Add("content", "false")
			}
		}

		// Handle array parameters
		if len(options.Lang) > 0 {
			params.Add("lang", strings.Join(options.Lang, ","))
		}
		if len(options.Country) > 0 {
			params.Add("country", strings.Join(options.Country, ","))
		}
		if len(options.Region) > 0 {
			params.Add("region", strings.Join(options.Region, ","))
		}
		if len(options.Category) > 0 {
			params.Add("category", strings.Join(options.Category, ","))
		}
		if len(options.Attributes) > 0 {
			params.Add("attributes", strings.Join(options.Attributes, ","))
		}
		if len(options.Publisher) > 0 {
			params.Add("publisher", strings.Join(options.Publisher, ","))
		}

		// Handle integer parameters
		if options.Max > 0 {
			params.Add("max", fmt.Sprintf("%d", options.Max))
		}
		if options.Page > 0 {
			params.Add("page", fmt.Sprintf("%d", options.Page))
		}

		// Handle other string parameters
		if options.SortBy != "" {
			params.Add("sortby", options.SortBy)
		}
		if options.Format != "" {
			params.Add("format", options.Format)
		}
	}

	// Build request URL
	headlinesURL := fmt.Sprintf("%s/v1/headlines?%s", c.baseURL, params.Encode())

	// Make the request
	req, err := http.NewRequest("GET", headlinesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Check for error responses
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, body)
	}

	// Parse the response
	var searchResponse SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &searchResponse, nil
}

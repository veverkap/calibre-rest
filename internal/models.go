package main

import (
	"fmt"
	"net/url"
	"strconv"
)

// Book represents a book in the Calibre library
type Book struct {
	Authors      []string          `json:"authors,omitempty"`
	AuthorSort   string            `json:"author_sort,omitempty"`
	Comments     string            `json:"comments,omitempty"`
	Cover        string            `json:"cover,omitempty"`
	Formats      []string          `json:"formats,omitempty"`
	ID           int               `json:"id,omitempty"`
	Identifiers  map[string]string `json:"identifiers,omitempty"`
	ISBN         string            `json:"isbn,omitempty"`
	Languages    []string          `json:"languages,omitempty"`
	LastModified string            `json:"last_modified,omitempty"`
	PubDate      string            `json:"pubdate,omitempty"`
	Publisher    string            `json:"publisher,omitempty"`
	Rating       int               `json:"rating,omitempty"`
	Series       string            `json:"series,omitempty"`
	SeriesIndex  float64           `json:"series_index,omitempty"`
	Size         int               `json:"size,omitempty"`
	Tags         []string          `json:"tags,omitempty"`
	Template     string            `json:"template,omitempty"`
	Timestamp    string            `json:"timestamp,omitempty"`
	Title        string            `json:"title,omitempty"`
	UUID         string            `json:"uuid,omitempty"`
}

// PaginatedResults represents paginated book results
type PaginatedResults struct {
	BaseURL *url.URL  `json:"-"`
	Books   []Book    `json:"-"`
	Start   int       `json:"-"`
	Limit   int       `json:"-"`
	Sort    []string  `json:"-"`
	Search  []string  `json:"-"`
	Count   int       `json:"-"`
}

// NewPaginatedResults creates a new paginated results object
func NewPaginatedResults(books []Book, start, limit int, sort, search []string) (*PaginatedResults, error) {
	baseURL, _ := url.Parse("/books")
	
	if len(books) < start {
		return nil, fmt.Errorf("start %d is larger than number of books (%d)", start, len(books))
	}
	
	return &PaginatedResults{
		BaseURL: baseURL,
		Books:   books,
		Start:   start,
		Limit:   limit,
		Sort:    sort,
		Search:  search,
		Count:   len(books),
	}, nil
}

// buildQuery builds a URL with query parameters
func (pr *PaginatedResults) buildQuery(start int) string {
	params := url.Values{}
	params.Set("start", strconv.Itoa(start))
	params.Set("limit", strconv.Itoa(pr.Limit))
	
	if pr.Sort != nil {
		for _, s := range pr.Sort {
			params.Add("sort", s)
		}
	}
	
	if pr.Search != nil {
		for _, s := range pr.Search {
			params.Add("search", s)
		}
	}
	
	u := *pr.BaseURL
	u.RawQuery = params.Encode()
	return u.String()
}

// CurrentPage returns the current page URL
func (pr *PaginatedResults) CurrentPage() string {
	return pr.buildQuery(pr.Start)
}

// PrevPage returns the previous page URL
func (pr *PaginatedResults) PrevPage() string {
	if !pr.HasPrevPage() {
		return ""
	}
	
	prevStart := pr.Start - pr.Limit
	if prevStart < 1 {
		prevStart = 1
	}
	
	return pr.buildQuery(prevStart)
}

// NextPage returns the next page URL
func (pr *PaginatedResults) NextPage() string {
	if !pr.HasNextPage() {
		return ""
	}
	
	return pr.buildQuery(pr.Start + pr.Limit)
}

// HasPrevPage returns true if there is a previous page
func (pr *PaginatedResults) HasPrevPage() bool {
	return pr.Start != 1
}

// HasNextPage returns true if there is a next page
func (pr *PaginatedResults) HasNextPage() bool {
	return pr.Start+pr.Limit <= pr.Count
}

// ToDict converts the paginated results to a dictionary-like structure
func (pr *PaginatedResults) ToDict() map[string]interface{} {
	endIndex := pr.Start - 1 + pr.Limit
	if endIndex > len(pr.Books) {
		endIndex = len(pr.Books)
	}
	
	books := make([]Book, 0)
	if pr.Start-1 < len(pr.Books) {
		books = pr.Books[pr.Start-1 : endIndex]
	}
	
	return map[string]interface{}{
		"books": books,
		"metadata": map[string]interface{}{
			"start": pr.Start,
			"limit": pr.Limit,
			"count": pr.Count,
			"self":  pr.CurrentPage(),
			"prev":  pr.PrevPage(),
			"next":  pr.NextPage(),
		},
	}
}
package db

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Core entity models representing database tables

// Author represents an author in the Calibre library
type Author struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Sort string `json:"sort" db:"sort"`
	Link string `json:"link" db:"link"`
}

// BookDB represents a book record as stored in the database
type BookDB struct {
	ID           int       `json:"id" db:"id"`
	Title        string    `json:"title" db:"title"`
	Sort         string    `json:"sort" db:"sort"`
	Timestamp    time.Time `json:"timestamp" db:"timestamp"`
	PubDate      time.Time `json:"pubdate" db:"pubdate"`
	SeriesIndex  float64   `json:"series_index" db:"series_index"`
	AuthorSort   string    `json:"author_sort" db:"author_sort"`
	ISBN         string    `json:"isbn" db:"isbn"`
	LCCN         string    `json:"lccn" db:"lccn"`
	Path         string    `json:"path" db:"path"`
	Flags        int       `json:"flags" db:"flags"`
	UUID         string    `json:"uuid" db:"uuid"`
	HasCover     bool      `json:"has_cover" db:"has_cover"`
	LastModified time.Time `json:"last_modified" db:"last_modified"`
}

// Series represents a book series
type Series struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Sort string `json:"sort" db:"sort"`
	Link string `json:"link" db:"link"`
}

// Tag represents a book tag
type Tag struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Link string `json:"link" db:"link"`
}

// Publisher represents a book publisher
type Publisher struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Sort string `json:"sort" db:"sort"`
	Link string `json:"link" db:"link"`
}

// Rating represents a book rating
type Rating struct {
	ID     int    `json:"id" db:"id"`
	Rating int    `json:"rating" db:"rating"`
	Link   string `json:"link" db:"link"`
}

// Language represents a language
type Language struct {
	ID       int    `json:"id" db:"id"`
	LangCode string `json:"lang_code" db:"lang_code"`
	Link     string `json:"link" db:"link"`
}

// Comment represents book comments/description
type Comment struct {
	ID   int    `json:"id" db:"id"`
	Book int    `json:"book" db:"book"`
	Text string `json:"text" db:"text"`
}

// Data represents book format data
type Data struct {
	ID               int    `json:"id" db:"id"`
	Book             int    `json:"book" db:"book"`
	Format           string `json:"format" db:"format"`
	UncompressedSize int    `json:"uncompressed_size" db:"uncompressed_size"`
	Name             string `json:"name" db:"name"`
}

// Identifier represents book identifiers (ISBN, etc.)
type Identifier struct {
	ID   int    `json:"id" db:"id"`
	Book int    `json:"book" db:"book"`
	Type string `json:"type" db:"type"`
	Val  string `json:"val" db:"val"`
}

// CustomColumn represents custom metadata columns
type CustomColumn struct {
	ID            int    `json:"id" db:"id"`
	Label         string `json:"label" db:"label"`
	Name          string `json:"name" db:"name"`
	Datatype      string `json:"datatype" db:"datatype"`
	MarkForDelete bool   `json:"mark_for_delete" db:"mark_for_delete"`
	Editable      bool   `json:"editable" db:"editable"`
	Display       string `json:"display" db:"display"`
	IsMultiple    bool   `json:"is_multiple" db:"is_multiple"`
	Normalized    bool   `json:"normalized" db:"normalized"`
}

// ConversionOption represents book conversion options
type ConversionOption struct {
	ID     int    `json:"id" db:"id"`
	Format string `json:"format" db:"format"`
	Book   *int   `json:"book" db:"book"` // nullable
	Data   []byte `json:"data" db:"data"`
}

// Feed represents RSS feeds
type Feed struct {
	ID     int    `json:"id" db:"id"`
	Title  string `json:"title" db:"title"`
	Script string `json:"script" db:"script"`
}

// LibraryID represents library identification
type LibraryID struct {
	ID   int    `json:"id" db:"id"`
	UUID string `json:"uuid" db:"uuid"`
}

// Preference represents user preferences
type Preference struct {
	ID  int    `json:"id" db:"id"`
	Key string `json:"key" db:"key"`
	Val string `json:"val" db:"val"`
}

// MetadataDirtied tracks books with dirty metadata
type MetadataDirtied struct {
	ID   int `json:"id" db:"id"`
	Book int `json:"book" db:"book"`
}

// AnnotationsDirtied tracks books with dirty annotations
type AnnotationsDirtied struct {
	ID   int `json:"id" db:"id"`
	Book int `json:"book" db:"book"`
}

// BooksPluginData represents plugin-specific book data
type BooksPluginData struct {
	ID   int    `json:"id" db:"id"`
	Book int    `json:"book" db:"book"`
	Name string `json:"name" db:"name"`
	Val  string `json:"val" db:"val"`
}

// LastReadPosition represents reading progress
type LastReadPosition struct {
	ID      int     `json:"id" db:"id"`
	Book    int     `json:"book" db:"book"`
	Format  string  `json:"format" db:"format"`
	User    string  `json:"user" db:"user"`
	Device  string  `json:"device" db:"device"`
	CFI     string  `json:"cfi" db:"cfi"`
	Epoch   float64 `json:"epoch" db:"epoch"`
	PosFrac float64 `json:"pos_frac" db:"pos_frac"`
}

// Annotation represents book annotations
type Annotation struct {
	ID             int     `json:"id" db:"id"`
	Book           int     `json:"book" db:"book"`
	Format         string  `json:"format" db:"format"`
	UserType       string  `json:"user_type" db:"user_type"`
	User           string  `json:"user" db:"user"`
	Timestamp      float64 `json:"timestamp" db:"timestamp"`
	AnnotID        string  `json:"annot_id" db:"annot_id"`
	AnnotType      string  `json:"annot_type" db:"annot_type"`
	AnnotData      string  `json:"annot_data" db:"annot_data"`
	SearchableText string  `json:"searchable_text" db:"searchable_text"`
}

// Link tables for many-to-many relationships

// BooksAuthorsLink represents the many-to-many relationship between books and authors
type BooksAuthorsLink struct {
	ID     int `json:"id" db:"id"`
	Book   int `json:"book" db:"book"`
	Author int `json:"author" db:"author"`
}

// BooksLanguagesLink represents the many-to-many relationship between books and languages
type BooksLanguagesLink struct {
	ID        int `json:"id" db:"id"`
	Book      int `json:"book" db:"book"`
	LangCode  int `json:"lang_code" db:"lang_code"`
	ItemOrder int `json:"item_order" db:"item_order"`
}

// BooksPublishersLink represents the many-to-many relationship between books and publishers
type BooksPublishersLink struct {
	ID        int `json:"id" db:"id"`
	Book      int `json:"book" db:"book"`
	Publisher int `json:"publisher" db:"publisher"`
}

// BooksRatingsLink represents the many-to-many relationship between books and ratings
type BooksRatingsLink struct {
	ID     int `json:"id" db:"id"`
	Book   int `json:"book" db:"book"`
	Rating int `json:"rating" db:"rating"`
}

// BooksSeriesLink represents the many-to-many relationship between books and series
type BooksSeriesLink struct {
	ID     int `json:"id" db:"id"`
	Book   int `json:"book" db:"book"`
	Series int `json:"series" db:"series"`
}

// BooksTagsLink represents the many-to-many relationship between books and tags
type BooksTagsLink struct {
	ID   int `json:"id" db:"id"`
	Book int `json:"book" db:"book"`
	Tag  int `json:"tag" db:"tag"`
}

// API response models

// Book represents a book in the Calibre library for API responses
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
	BaseURL *url.URL `json:"-"`
	Books   []Book   `json:"-"`
	Start   int      `json:"-"`
	Limit   int      `json:"-"`
	Sort    []string `json:"-"`
	Search  []string `json:"-"`
	Count   int      `json:"-"`
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

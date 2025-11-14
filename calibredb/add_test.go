package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_AddHelp(t *testing.T) {
	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()
	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
	)

	got := c.AddHelp()
	if got == "" {
		t.Error("AddHelp() returned empty string")
	}
}

func TestCalibre_Add(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		opts    calibredb.AddOptions
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "Missing required Files field",
			opts: calibredb.AddOptions{
				Files: []string{},
			},
			wantErr: true,
		},
		{
			name: "Valid with single file",
			opts: calibredb.AddOptions{
				Files: []string{"book1.epub"},
			},
			wantErr: false,
		},
		{
			name: "Valid with multiple files",
			opts: calibredb.AddOptions{
				Files: []string{"book1.epub", "book2.pdf"},
			},
			wantErr: false,
		},
		{
			name: "With Authors option",
			opts: calibredb.AddOptions{
				Files:   []string{"book.epub"},
				Authors: "John Doe",
			},
			wantErr: false,
		},
		{
			name: "With Automerge disabled",
			opts: calibredb.AddOptions{
				Files:     []string{"book.epub"},
				Automerge: calibredb.Disabled,
			},
			wantErr: false,
		},
		{
			name: "With Automerge ignore",
			opts: calibredb.AddOptions{
				Files:     []string{"book.epub"},
				Automerge: calibredb.Ignore,
			},
			wantErr: false,
		},
		{
			name: "With Automerge overwrite",
			opts: calibredb.AddOptions{
				Files:     []string{"book.epub"},
				Automerge: calibredb.Overwrite,
			},
			wantErr: false,
		},
		{
			name: "With Automerge new_record",
			opts: calibredb.AddOptions{
				Files:     []string{"book.epub"},
				Automerge: calibredb.NewRecord,
			},
			wantErr: false,
		},
		{
			name: "With Cover option",
			opts: calibredb.AddOptions{
				Files: []string{"book.epub"},
				Cover: "/path/to/cover.jpg",
			},
			wantErr: false,
		},
		{
			name: "With Duplicates true",
			opts: calibredb.AddOptions{
				Files:      []string{"book.epub"},
				Duplicates: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With Duplicates false",
			opts: calibredb.AddOptions{
				Files:      []string{"book.epub"},
				Duplicates: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "With Empty true",
			opts: calibredb.AddOptions{
				Files: []string{"book.epub"},
				Empty: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With Empty false",
			opts: calibredb.AddOptions{
				Files: []string{"book.epub"},
				Empty: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "With single Identifier",
			opts: calibredb.AddOptions{
				Files:      []string{"book.epub"},
				Identifier: []string{"isbn:1234567890"},
			},
			wantErr: false,
		},
		{
			name: "With multiple Identifiers",
			opts: calibredb.AddOptions{
				Files:      []string{"book.epub"},
				Identifier: []string{"isbn:1234567890", "asin:ABCD123"},
			},
			wantErr: false,
		},
		{
			name: "With ISBN option",
			opts: calibredb.AddOptions{
				Files: []string{"book.epub"},
				Isbn:  "1234567890",
			},
			wantErr: false,
		},
		{
			name: "With Languages option",
			opts: calibredb.AddOptions{
				Files:     []string{"book.epub"},
				Languages: "eng,fra",
			},
			wantErr: false,
		},
		{
			name: "With Series option",
			opts: calibredb.AddOptions{
				Files:  []string{"book.epub"},
				Series: "The Great Series",
			},
			wantErr: false,
		},
		{
			name: "With SeriesIndex option",
			opts: calibredb.AddOptions{
				Files:       []string{"book.epub"},
				SeriesIndex: 3.5,
			},
			wantErr: false,
		},
		{
			name: "With SeriesIndex zero (should be omitted)",
			opts: calibredb.AddOptions{
				Files:       []string{"book.epub"},
				SeriesIndex: 0,
			},
			wantErr: false,
		},
		{
			name: "With Tags option",
			opts: calibredb.AddOptions{
				Files: []string{"book.epub"},
				Tags:  "fiction,adventure",
			},
			wantErr: false,
		},
		{
			name: "With Title option",
			opts: calibredb.AddOptions{
				Files: []string{"book.epub"},
				Title: "My Great Book",
			},
			wantErr: false,
		},
		{
			name: "With OneBookPerDirectory true",
			opts: calibredb.AddOptions{
				Files:               []string{"book.epub"},
				OneBookPerDirectory: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With OneBookPerDirectory false",
			opts: calibredb.AddOptions{
				Files:               []string{"book.epub"},
				OneBookPerDirectory: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "With Recurse true",
			opts: calibredb.AddOptions{
				Files:   []string{"book.epub"},
				Recurse: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With Recurse false",
			opts: calibredb.AddOptions{
				Files:   []string{"book.epub"},
				Recurse: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "With all options combined",
			opts: calibredb.AddOptions{
				Files:               []string{"book1.epub", "book2.pdf"},
				Authors:             "Jane Smith",
				Automerge:           calibredb.Overwrite,
				Cover:               "/path/to/cover.jpg",
				Duplicates:          func(b bool) *bool { return &b }(true),
				Empty:               func(b bool) *bool { return &b }(false),
				Identifier:          []string{"isbn:1234567890", "asin:ABCD123"},
				Isbn:                "1234567890",
				Languages:           "eng,spa",
				Series:              "Epic Series",
				SeriesIndex:         2.5,
				Tags:                "fiction,scifi",
				Title:               "Epic Book Title",
				OneBookPerDirectory: func(b bool) *bool { return &b }(true),
				Recurse:             func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := os.TempDir() + "/" + t.Name()
			defer func() { _ = os.RemoveAll(tempDir) }()
			c := calibredb.NewCalibre(
				calibredb.WithLibraryPath(tempDir),
				calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
			)

			got, gotErr := c.Add(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					// Check if the error is just because calibredb is not installed
					if strings.Contains(gotErr.Error(), "no such file or directory") {
						t.Skip("Skipping test: calibredb not found")
					}
					t.Errorf("Add() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Add() succeeded unexpectedly")
			}
			if tt.want != "" && !strings.HasPrefix(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAutomergeChoice_Constants(t *testing.T) {
	tests := []struct {
		name     string
		choice   calibredb.AutomergeChoice
		expected string
	}{
		{
			name:     "Disabled constant",
			choice:   calibredb.Disabled,
			expected: "disabled",
		},
		{
			name:     "Ignore constant",
			choice:   calibredb.Ignore,
			expected: "ignore",
		},
		{
			name:     "Overwrite constant",
			choice:   calibredb.Overwrite,
			expected: "overwrite",
		},
		{
			name:     "NewRecord constant",
			choice:   calibredb.NewRecord,
			expected: "new_record",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.choice) != tt.expected {
				t.Errorf("AutomergeChoice constant = %v, want %v", tt.choice, tt.expected)
			}
		})
	}
}

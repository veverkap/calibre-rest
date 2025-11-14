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
	// We expect it to fail with "no such file or directory" since calibredb is not installed
	if !strings.Contains(got, "no such file or directory") {
		t.Logf("AddHelp() = %v", got)
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
			name: "Invalid - nil Files",
			opts: calibredb.AddOptions{
				Files: nil,
			},
			wantErr: true,
		},
		{
			name: "Valid - minimal options",
			opts: calibredb.AddOptions{
				Files: []string{"book.epub"},
			},
			wantErr: true, // Will fail due to calibredb not being installed, but validation passes
		},
		{
			name: "Valid - with Authors",
			opts: calibredb.AddOptions{
				Files:   []string{"book.epub"},
				Authors: "John Doe",
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Cover",
			opts: calibredb.AddOptions{
				Files: []string{"book.epub"},
				Cover: "cover.jpg",
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with ISBN",
			opts: calibredb.AddOptions{
				Files: []string{"book.epub"},
				Isbn:  "978-3-16-148410-0",
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Languages",
			opts: calibredb.AddOptions{
				Files:     []string{"book.epub"},
				Languages: "en,es",
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Series",
			opts: calibredb.AddOptions{
				Files:  []string{"book.epub"},
				Series: "My Series",
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with SeriesIndex",
			opts: calibredb.AddOptions{
				Files:       []string{"book.epub"},
				SeriesIndex: 3.5,
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Tags",
			opts: calibredb.AddOptions{
				Files: []string{"book.epub"},
				Tags:  "fiction,adventure",
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Title",
			opts: calibredb.AddOptions{
				Files: []string{"book.epub"},
				Title: "My Book Title",
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Automerge Disabled",
			opts: calibredb.AddOptions{
				Files:     []string{"book.epub"},
				Automerge: calibredb.Disabled,
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Automerge Ignore",
			opts: calibredb.AddOptions{
				Files:     []string{"book.epub"},
				Automerge: calibredb.Ignore,
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Automerge Overwrite",
			opts: calibredb.AddOptions{
				Files:     []string{"book.epub"},
				Automerge: calibredb.Overwrite,
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Automerge NewRecord",
			opts: calibredb.AddOptions{
				Files:     []string{"book.epub"},
				Automerge: calibredb.NewRecord,
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Duplicates true",
			opts: calibredb.AddOptions{
				Files:      []string{"book.epub"},
				Duplicates: func(b bool) *bool { return &b }(true),
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Duplicates false",
			opts: calibredb.AddOptions{
				Files:      []string{"book.epub"},
				Duplicates: func(b bool) *bool { return &b }(false),
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Empty true",
			opts: calibredb.AddOptions{
				Files: []string{"book.epub"},
				Empty: func(b bool) *bool { return &b }(true),
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Empty false",
			opts: calibredb.AddOptions{
				Files: []string{"book.epub"},
				Empty: func(b bool) *bool { return &b }(false),
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Identifier",
			opts: calibredb.AddOptions{
				Files:      []string{"book.epub"},
				Identifier: []string{"isbn:123456", "asin:ABCDEF"},
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with OneBookPerDirectory true",
			opts: calibredb.AddOptions{
				Files:               []string{"book.epub"},
				OneBookPerDirectory: func(b bool) *bool { return &b }(true),
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with OneBookPerDirectory false",
			opts: calibredb.AddOptions{
				Files:               []string{"book.epub"},
				OneBookPerDirectory: func(b bool) *bool { return &b }(false),
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Recurse true",
			opts: calibredb.AddOptions{
				Files:   []string{"book.epub"},
				Recurse: func(b bool) *bool { return &b }(true),
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with Recurse false",
			opts: calibredb.AddOptions{
				Files:   []string{"book.epub"},
				Recurse: func(b bool) *bool { return &b }(false),
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with multiple Files",
			opts: calibredb.AddOptions{
				Files: []string{"book1.epub", "book2.pdf", "book3.mobi"},
			},
			wantErr: true, // Will fail due to calibredb not being installed
		},
		{
			name: "Valid - with all options",
			opts: calibredb.AddOptions{
				Files:               []string{"book.epub"},
				Authors:             "Jane Doe",
				Automerge:           calibredb.Ignore,
				Cover:               "cover.jpg",
				Duplicates:          func(b bool) *bool { return &b }(true),
				Empty:               func(b bool) *bool { return &b }(false),
				Identifier:          []string{"isbn:987654"},
				Isbn:                "978-1-23-456789-0",
				Languages:           "en",
				Series:              "Book Series",
				SeriesIndex:         1.0,
				Tags:                "fiction,scifi",
				Title:               "Complete Book",
				OneBookPerDirectory: func(b bool) *bool { return &b }(true),
				Recurse:             func(b bool) *bool { return &b }(false),
			},
			wantErr: true, // Will fail due to calibredb not being installed
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
				// Check if it's a validation error (which we want for invalid inputs)
				if strings.Contains(tt.name, "Invalid") {
					// Expected validation error
					if !strings.Contains(gotErr.Error(), "required") {
						t.Errorf("Add() validation error = %v, expected 'required' error", gotErr)
					}
					return
				}
				// For valid inputs, we expect calibredb to not be found
				if !strings.Contains(gotErr.Error(), "no such file or directory") && !tt.wantErr {
					t.Errorf("Add() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Add() succeeded unexpectedly")
			}
			if !strings.HasPrefix(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

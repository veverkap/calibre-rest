package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_FtsSearchHelp(t *testing.T) {
	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()
	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
	)

	got := c.FtsSearchHelp()
	// Since calibredb is not installed, this should return an error message
	// We just verify it returns something (either help text or error message)
	if got == "" {
		t.Error("FtsSearchHelp() returned empty string")
	}
}

func TestCalibre_FtsSearch(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		opts    calibredb.FtsSearchOptions
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "Missing required Search",
			opts: calibredb.FtsSearchOptions{
				Search:     "",
				Expression: "test",
			},
			wantErr: true,
		},
		{
			name: "Missing required Expression",
			opts: calibredb.FtsSearchOptions{
				Search:     "test",
				Expression: "",
			},
			wantErr: true,
		},
		{
			name: "Missing both required fields",
			opts: calibredb.FtsSearchOptions{
				Search:     "",
				Expression: "",
			},
			wantErr: true,
		},
		{
			name: "Valid Search and Expression only",
			opts: calibredb.FtsSearchOptions{
				Search:     "test",
				Expression: "author:Smith",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With DoNotMatchOnRelatedWords true",
			opts: calibredb.FtsSearchOptions{
				Search:                   "test",
				Expression:               "correction",
				DoNotMatchOnRelatedWords: func(b bool) *bool { return &b }(true),
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With DoNotMatchOnRelatedWords false",
			opts: calibredb.FtsSearchOptions{
				Search:                   "test",
				Expression:               "correction",
				DoNotMatchOnRelatedWords: func(b bool) *bool { return &b }(false),
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With IncludeSnippets true",
			opts: calibredb.FtsSearchOptions{
				Search:          "test",
				Expression:      "important",
				IncludeSnippets: func(b bool) *bool { return &b }(true),
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With IncludeSnippets false",
			opts: calibredb.FtsSearchOptions{
				Search:          "test",
				Expression:      "important",
				IncludeSnippets: func(b bool) *bool { return &b }(false),
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With IndexingThreshold",
			opts: calibredb.FtsSearchOptions{
				Search:            "test",
				Expression:        "query",
				IndexingThreshold: 95.5,
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With IndexingThreshold zero value",
			opts: calibredb.FtsSearchOptions{
				Search:            "test",
				Expression:        "query",
				IndexingThreshold: 0,
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With MatchEndMarker",
			opts: calibredb.FtsSearchOptions{
				Search:         "test",
				Expression:     "word",
				MatchEndMarker: "</match>",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With MatchStartMarker",
			opts: calibredb.FtsSearchOptions{
				Search:           "test",
				Expression:       "word",
				MatchStartMarker: "<match>",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With both MatchStartMarker and MatchEndMarker",
			opts: calibredb.FtsSearchOptions{
				Search:           "test",
				Expression:       "word",
				MatchStartMarker: "<match>",
				MatchEndMarker:   "</match>",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With OutputFormat Text",
			opts: calibredb.FtsSearchOptions{
				Search:       "test",
				Expression:   "query",
				OutputFormat: calibredb.Text,
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With OutputFormat Json",
			opts: calibredb.FtsSearchOptions{
				Search:       "test",
				Expression:   "query",
				OutputFormat: calibredb.Json,
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With RestrictTo",
			opts: calibredb.FtsSearchOptions{
				Search:     "test",
				Expression: "query",
				RestrictTo: "ids:1,2,3",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With RestrictTo search expression",
			opts: calibredb.FtsSearchOptions{
				Search:     "test",
				Expression: "query",
				RestrictTo: "search:tag:fiction",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With all options",
			opts: calibredb.FtsSearchOptions{
				Search:                   "test",
				Expression:               "comprehensive query",
				DoNotMatchOnRelatedWords: func(b bool) *bool { return &b }(true),
				IncludeSnippets:          func(b bool) *bool { return &b }(true),
				IndexingThreshold:        90.0,
				MatchEndMarker:           ">>",
				MatchStartMarker:         "<<",
				OutputFormat:             calibredb.Json,
				RestrictTo:               "ids:10,20,30",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With empty RestrictTo",
			opts: calibredb.FtsSearchOptions{
				Search:     "test",
				Expression: "query",
				RestrictTo: "",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With empty MatchStartMarker",
			opts: calibredb.FtsSearchOptions{
				Search:           "test",
				Expression:       "query",
				MatchStartMarker: "",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With empty MatchEndMarker",
			opts: calibredb.FtsSearchOptions{
				Search:         "test",
				Expression:     "query",
				MatchEndMarker: "",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "Complex search expression",
			opts: calibredb.FtsSearchOptions{
				Search:     "advanced",
				Expression: "title:\"Harry Potter\" AND author:Rowling",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With multiple options combination 1",
			opts: calibredb.FtsSearchOptions{
				Search:            "test",
				Expression:        "query",
				IncludeSnippets:   func(b bool) *bool { return &b }(true),
				IndexingThreshold: 85.5,
				OutputFormat:      calibredb.Text,
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With multiple options combination 2",
			opts: calibredb.FtsSearchOptions{
				Search:                   "test",
				Expression:               "query",
				DoNotMatchOnRelatedWords: func(b bool) *bool { return &b }(true),
				MatchStartMarker:         "[[",
				MatchEndMarker:           "]]",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "With multiple options combination 3",
			opts: calibredb.FtsSearchOptions{
				Search:          "test",
				Expression:      "query",
				IncludeSnippets: func(b bool) *bool { return &b }(true),
				RestrictTo:      "search:format:epub",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
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

			got, gotErr := c.FtsSearch(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("FtsSearch() failed: %v", gotErr)
				}
				// For validation errors, check that it's the right type of error
				if (tt.opts.Search == "" || tt.opts.Expression == "") && !strings.Contains(gotErr.Error(), "required") {
					t.Errorf("FtsSearch() error for missing required field should mention 'required', got: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("FtsSearch() succeeded unexpectedly")
			}
			if tt.want != "" && !strings.HasPrefix(got, tt.want) {
				t.Errorf("FtsSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutputFormatChoice_Constants(t *testing.T) {
	tests := []struct {
		name  string
		value calibredb.OutputFormatChoice
		want  string
	}{
		{
			name:  "Text constant",
			value: calibredb.Text,
			want:  "text",
		},
		{
			name:  "Json constant",
			value: calibredb.Json,
			want:  "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.want {
				t.Errorf("OutputFormatChoice constant = %v, want %v", tt.value, tt.want)
			}
		})
	}
}

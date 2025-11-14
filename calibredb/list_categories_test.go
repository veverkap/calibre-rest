package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_ListCategoriesHelp(t *testing.T) {
	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()
	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
	)

	got := c.ListCategoriesHelp()
	if got == "" {
		t.Error("ListCategoriesHelp() returned empty string")
	}
}

func TestCalibre_ListCategories(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		opts    calibredb.ListCategoriesOptions
		args    []string
		want    string
		wantErr bool
	}{
		{
			name:    "Empty options - should succeed",
			opts:    calibredb.ListCategoriesOptions{},
			wantErr: false,
		},
		{
			name: "With Categories single",
			opts: calibredb.ListCategoriesOptions{
				Categories: "tags",
			},
			wantErr: false,
		},
		{
			name: "With Categories multiple",
			opts: calibredb.ListCategoriesOptions{
				Categories: "tags,authors,series",
			},
			wantErr: false,
		},
		{
			name: "With Categories all common types",
			opts: calibredb.ListCategoriesOptions{
				Categories: "tags,authors,series,publisher,languages",
			},
			wantErr: false,
		},
		{
			name: "With Csv true",
			opts: calibredb.ListCategoriesOptions{
				Csv: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With Csv false",
			opts: calibredb.ListCategoriesOptions{
				Csv: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "With Dialect excel",
			opts: calibredb.ListCategoriesOptions{
				Dialect: calibredb.DialectExcel,
			},
			wantErr: false,
		},
		{
			name: "With Dialect excel-tab",
			opts: calibredb.ListCategoriesOptions{
				Dialect: calibredb.DialectExcelTab,
			},
			wantErr: false,
		},
		{
			name: "With Dialect unix",
			opts: calibredb.ListCategoriesOptions{
				Dialect: calibredb.DialectUnix,
			},
			wantErr: false,
		},
		{
			name: "With ItemCount true",
			opts: calibredb.ListCategoriesOptions{
				ItemCount: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With ItemCount false",
			opts: calibredb.ListCategoriesOptions{
				ItemCount: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "With Width 80",
			opts: calibredb.ListCategoriesOptions{
				Width: 80,
			},
			wantErr: false,
		},
		{
			name: "With Width 120",
			opts: calibredb.ListCategoriesOptions{
				Width: 120,
			},
			wantErr: false,
		},
		{
			name: "With Width 0 (should be omitted)",
			opts: calibredb.ListCategoriesOptions{
				Width: 0,
			},
			wantErr: false,
		},
		{
			name: "With Categories and Csv",
			opts: calibredb.ListCategoriesOptions{
				Categories: "tags,authors",
				Csv:        func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With Csv and Dialect",
			opts: calibredb.ListCategoriesOptions{
				Csv:     func(b bool) *bool { return &b }(true),
				Dialect: calibredb.DialectUnix,
			},
			wantErr: false,
		},
		{
			name: "With Csv and ItemCount",
			opts: calibredb.ListCategoriesOptions{
				Csv:       func(b bool) *bool { return &b }(true),
				ItemCount: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With Categories and Width",
			opts: calibredb.ListCategoriesOptions{
				Categories: "series",
				Width:      100,
			},
			wantErr: false,
		},
		{
			name: "With all options combined",
			opts: calibredb.ListCategoriesOptions{
				Categories: "tags,authors,series,publisher",
				Csv:        func(b bool) *bool { return &b }(true),
				Dialect:    calibredb.DialectExcel,
				ItemCount:  func(b bool) *bool { return &b }(true),
				Width:      120,
			},
			wantErr: false,
		},
		{
			name: "With all options combined (different dialect)",
			opts: calibredb.ListCategoriesOptions{
				Categories: "languages,formats",
				Csv:        func(b bool) *bool { return &b }(true),
				Dialect:    calibredb.DialectExcelTab,
				ItemCount:  func(b bool) *bool { return &b }(false),
				Width:      80,
			},
			wantErr: false,
		},
		{
			name: "With all options combined (unix dialect)",
			opts: calibredb.ListCategoriesOptions{
				Categories: "authors",
				Csv:        func(b bool) *bool { return &b }(true),
				Dialect:    calibredb.DialectUnix,
				ItemCount:  func(b bool) *bool { return &b }(true),
				Width:      160,
			},
			wantErr: false,
		},
		{
			name: "With empty Categories string",
			opts: calibredb.ListCategoriesOptions{
				Categories: "",
			},
			wantErr: false,
		},
		{
			name: "With ItemCount only",
			opts: calibredb.ListCategoriesOptions{
				ItemCount: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With Dialect without Csv",
			opts: calibredb.ListCategoriesOptions{
				Dialect: calibredb.DialectExcel,
			},
			wantErr: false,
		},
		{
			name: "With Width and ItemCount",
			opts: calibredb.ListCategoriesOptions{
				Width:     100,
				ItemCount: func(b bool) *bool { return &b }(true),
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

			got, gotErr := c.ListCategories(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					// Check if the error is just because calibredb is not installed
					if strings.Contains(gotErr.Error(), "no such file or directory") {
						t.Skip("Skipping test: calibredb not found")
					}
					t.Errorf("ListCategories() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ListCategories() succeeded unexpectedly")
			}
			if tt.want != "" && !strings.HasPrefix(got, tt.want) {
				t.Errorf("ListCategories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDialectChoice_Constants(t *testing.T) {
	tests := []struct {
		name     string
		choice   calibredb.DialectChoice
		expected string
	}{
		{
			name:     "DialectExcel constant",
			choice:   calibredb.DialectExcel,
			expected: "excel",
		},
		{
			name:     "DialectExcelTab constant",
			choice:   calibredb.DialectExcelTab,
			expected: "excel-tab",
		},
		{
			name:     "DialectUnix constant",
			choice:   calibredb.DialectUnix,
			expected: "unix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.choice) != tt.expected {
				t.Errorf("DialectChoice constant = %v, want %v", tt.choice, tt.expected)
			}
		})
	}
}

package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_SetMetadataHelp(t *testing.T) {
	c, f := getTestCalibre(t.Name())
	defer f()
	got := c.SetMetadataHelp()
	if got == "" {
		t.Error("SetMetadataHelp() returned empty string")
	}
}

func TestCalibre_SetMetadata(t *testing.T) {
	tests := []struct {
		name    string
		opts    calibredb.SetMetadataOptions
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "Missing required BookId field",
			opts: calibredb.SetMetadataOptions{
				BookId: "",
				Path:   "/path/to/metadata.opf",
			},
			wantErr: true,
		},
		{
			name: "Missing required Path field",
			opts: calibredb.SetMetadataOptions{
				BookId: "1",
				Path:   "",
			},
			wantErr: true,
		},
		{
			name: "Both required fields missing",
			opts: calibredb.SetMetadataOptions{
				BookId: "",
				Path:   "",
			},
			wantErr: true,
		},
		{
			name: "Valid with minimal options",
			opts: calibredb.SetMetadataOptions{
				BookId: "1",
				Path:   "/path/to/metadata.opf",
			},
			wantErr: false,
		},
		{
			name: "With single Field option",
			opts: calibredb.SetMetadataOptions{
				BookId: "1",
				Path:   "/path/to/metadata.opf",
				Field:  []string{"title:My Book Title"},
			},
			wantErr: false,
		},
		{
			name: "With multiple Field options",
			opts: calibredb.SetMetadataOptions{
				BookId: "1",
				Path:   "/path/to/metadata.opf",
				Field: []string{
					"title:My Book Title",
					"authors:John Doe",
					"tags:fiction,adventure",
				},
			},
			wantErr: false,
		},
		{
			name: "With Field option for languages",
			opts: calibredb.SetMetadataOptions{
				BookId: "1",
				Path:   "/path/to/metadata.opf",
				Field:  []string{"languages:en,fr"},
			},
			wantErr: false,
		},
		{
			name: "With Field option for identifiers",
			opts: calibredb.SetMetadataOptions{
				BookId: "1",
				Path:   "/path/to/metadata.opf",
				Field:  []string{"identifiers:isbn:1234567890,doi:10.1234/example"},
			},
			wantErr: false,
		},
		{
			name: "With Field option for boolean field",
			opts: calibredb.SetMetadataOptions{
				BookId: "1",
				Path:   "/path/to/metadata.opf",
				Field:  []string{"some_boolean:true"},
			},
			wantErr: false,
		},
		{
			name: "With ListFields set to true",
			opts: calibredb.SetMetadataOptions{
				BookId:     "1",
				Path:       "/path/to/metadata.opf",
				ListFields: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With ListFields set to false",
			opts: calibredb.SetMetadataOptions{
				BookId:     "1",
				Path:       "/path/to/metadata.opf",
				ListFields: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "With ListFields nil (should not add flag)",
			opts: calibredb.SetMetadataOptions{
				BookId:     "1",
				Path:       "/path/to/metadata.opf",
				ListFields: nil,
			},
			wantErr: false,
		},
		{
			name: "With empty Field slice (should not add flag)",
			opts: calibredb.SetMetadataOptions{
				BookId: "1",
				Path:   "/path/to/metadata.opf",
				Field:  []string{},
			},
			wantErr: false,
		},
		{
			name: "With all options combined",
			opts: calibredb.SetMetadataOptions{
				BookId: "42",
				Path:   "/path/to/book/metadata.opf",
				Field: []string{
					"title:Complete Book Title",
					"authors:Jane Smith,John Doe",
					"tags:scifi,adventure,bestseller",
					"languages:en",
					"identifiers:isbn:9876543210,asin:B0EXAMPLE",
					"series:Epic Series",
					"series_index:3",
				},
				ListFields: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "With numeric BookId",
			opts: calibredb.SetMetadataOptions{
				BookId: "999",
				Path:   "/path/to/metadata.opf",
			},
			wantErr: false,
		},
		{
			name: "With Field for series and series_index",
			opts: calibredb.SetMetadataOptions{
				BookId: "1",
				Path:   "/path/to/metadata.opf",
				Field: []string{
					"series:The Great Series",
					"series_index:2.5",
				},
			},
			wantErr: false,
		},
		{
			name: "With Field for publisher and pubdate",
			opts: calibredb.SetMetadataOptions{
				BookId: "1",
				Path:   "/path/to/metadata.opf",
				Field: []string{
					"publisher:Great Publisher",
					"pubdate:2023-01-15",
				},
			},
			wantErr: false,
		},
		{
			name: "With Field for rating and comments",
			opts: calibredb.SetMetadataOptions{
				BookId: "1",
				Path:   "/path/to/metadata.opf",
				Field: []string{
					"rating:5",
					"comments:This is a great book with a long description.",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, f := getTestCalibre(t.Name())
			defer f()

			got, gotErr := c.SetMetadata(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					// Check if the error is just because calibredb is not installed
					if strings.Contains(gotErr.Error(), "no such file or directory") ||
						strings.Contains(gotErr.Error(), "executable file not found") {
						t.Skip("Skipping test: calibredb not found")
					}
					t.Errorf("SetMetadata() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SetMetadata() succeeded unexpectedly")
			}
			if tt.want != "" && !strings.HasPrefix(got, tt.want) {
				t.Errorf("SetMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

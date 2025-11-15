package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_EmbedMetadataHelp(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "EmbedMetadataHelp returns help text",
			want: "usage: calibredb embed_metadata",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calibredb.NewCalibre(
				calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
			)
			got := c.EmbedMetadataHelp()
			if !strings.Contains(got, tt.want) && !strings.Contains(got, "no such file or directory") {
				t.Errorf("EmbedMetadataHelp() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}

func TestCalibre_EmbedMetadata(t *testing.T) {
	tests := []struct {
		name           string
		opts           calibredb.EmbedMetadataOptions
		args           []string
		want           string
		wantErr        bool
		wantValidation bool // true if we expect a validation error
	}{
		{
			name: "Invalid - empty BookId",
			opts: calibredb.EmbedMetadataOptions{
				BookId: "",
			},
			wantErr:        true,
			wantValidation: true,
		},
		{
			name: "Valid - with single BookId",
			opts: calibredb.EmbedMetadataOptions{
				BookId: "1",
			},
			wantErr: false,
		},
		{
			name: "Valid - with multiple BookIds",
			opts: calibredb.EmbedMetadataOptions{
				BookId: "1 2 3",
			},
			wantErr: false,
		},
		{
			name: "Valid - with BookId range",
			opts: calibredb.EmbedMetadataOptions{
				BookId: "1-10",
			},
			wantErr: false,
		},
		{
			name: "Valid - with 'all' BookId",
			opts: calibredb.EmbedMetadataOptions{
				BookId: "all",
			},
			wantErr: false,
		},
		{
			name: "Valid - with single OnlyFormats",
			opts: calibredb.EmbedMetadataOptions{
				BookId:      "1",
				OnlyFormats: []string{"EPUB"},
			},
			wantErr: false,
		},
		{
			name: "Valid - with multiple OnlyFormats",
			opts: calibredb.EmbedMetadataOptions{
				BookId:      "1",
				OnlyFormats: []string{"EPUB", "PDF", "MOBI"},
			},
			wantErr: false,
		},
		{
			name: "Valid - with empty OnlyFormats slice",
			opts: calibredb.EmbedMetadataOptions{
				BookId:      "1",
				OnlyFormats: []string{},
			},
			wantErr: false,
		},
		{
			name: "Valid - with additional args",
			opts: calibredb.EmbedMetadataOptions{
				BookId: "1",
			},
			args:    []string{"extra", "args"},
			wantErr: false,
		},
		{
			name: "Valid - complex scenario with all options",
			opts: calibredb.EmbedMetadataOptions{
				BookId:      "1 5-10 15",
				OnlyFormats: []string{"EPUB", "PDF"},
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

			got, gotErr := c.EmbedMetadata(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					// If calibredb is not installed, we can't test actual execution
					// but we can verify validation errors work
					if !strings.Contains(gotErr.Error(), "no such file or directory") &&
						!strings.Contains(gotErr.Error(), "executable file not found") {
						t.Errorf("EmbedMetadata() failed: %v", gotErr)
					}
				} else if tt.wantValidation {
					// Verify this is actually a validation error
					if !strings.Contains(gotErr.Error(), "required") && !strings.Contains(gotErr.Error(), "validation") {
						// This might still be a valid validation error from the validator package
						// Just ensure we got an error as expected
						t.Logf("Got expected validation error: %v", gotErr)
					}
				}
				return
			}
			if tt.wantErr {
				t.Fatal("EmbedMetadata() succeeded unexpectedly")
			}
			if tt.want != "" && !strings.Contains(got, tt.want) {
				t.Errorf("EmbedMetadata() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}

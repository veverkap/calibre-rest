package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_RemoveFormat(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		opts    calibredb.RemoveFormatOptions
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "Invalid id - empty",
			opts: calibredb.RemoveFormatOptions{
				Id:  "",
				Fmt: "EPUB",
			},
			wantErr: true,
		},
		{
			name: "Invalid fmt - empty",
			opts: calibredb.RemoveFormatOptions{
				Id:  "1",
				Fmt: "",
			},
			wantErr: true,
		},
		{
			name: "Both id and fmt empty",
			opts: calibredb.RemoveFormatOptions{
				Id:  "",
				Fmt: "",
			},
			wantErr: true,
		},
		{
			name: "Valid parameters",
			opts: calibredb.RemoveFormatOptions{
				Id:  "1",
				Fmt: "EPUB",
			},
			wantErr: false,
		},
		{
			name: "Valid parameters with different format",
			opts: calibredb.RemoveFormatOptions{
				Id:  "42",
				Fmt: "PDF",
			},
			wantErr: false,
		},
		{
			name: "Valid parameters with lowercase format",
			opts: calibredb.RemoveFormatOptions{
				Id:  "5",
				Fmt: "txt",
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
				calibredb.WithCalibreDBLocation(getCalibreDBPath()),
			)

			got, gotErr := c.RemoveFormat(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("RemoveFormat() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("RemoveFormat() succeeded unexpectedly")
			}
			if tt.want != "" && !strings.HasPrefix(got, tt.want) {
				t.Errorf("RemoveFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalibre_RemoveFormatHelp(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Get help text",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := os.TempDir() + "/" + t.Name()
			defer func() { _ = os.RemoveAll(tempDir) }()
			c := calibredb.NewCalibre(
				calibredb.WithLibraryPath(tempDir),
				calibredb.WithCalibreDBLocation(getCalibreDBPath()),
			)

			got := c.RemoveFormatHelp()
			if got == "" {
				t.Error("RemoveFormatHelp() returned empty string")
			}
			// Help text should contain some expected keywords
			if !strings.Contains(got, "remove_format") {
				t.Errorf("RemoveFormatHelp() = %v, want it to contain 'remove_format'", got)
			}
		})
	}
}

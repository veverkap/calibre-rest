package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_AddFormat(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		opts    calibredb.AddFormatOptions
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "Invalid id",
			opts: calibredb.AddFormatOptions{
				Id:        "",
				EbookFile: "somefile.epub",
			},
			wantErr: true,
		},
		{
			name: "Invalid ebook file",

			opts: calibredb.AddFormatOptions{
				Id:        "someid",
				EbookFile: "",
			},
			wantErr: true,
		},
		// {
		// 	name: "Valid parameters",
		// 	opts: calibredb.AddFormatOptions{
		// 		Id:        "1",
		// 		EbookFile: "somefile.epub",
		// 	},
		// 	want:    "Added format",
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := os.WriteFile("somefile.epub", []byte("ok"), 0644)
			if err != nil {
				t.Fatalf("Failed to create temporary ebook file: %v", err)
			}
			defer func() { _ = os.Remove("somefile.epub") }()
			tempDir := os.TempDir() + "/" + t.Name()
			defer func() { _ = os.RemoveAll(tempDir) }()
			c := calibredb.NewCalibre(
				calibredb.WithLibraryPath(tempDir),
				calibredb.WithCalibreDBLocation(getCalibreDBPath()),
			)

			got, gotErr := c.AddFormat(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("AddFormat() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("AddFormat() succeeded unexpectedly")
			}
			if !strings.HasPrefix(got, tt.want) {
				t.Errorf("AddFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

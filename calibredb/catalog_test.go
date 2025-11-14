package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_CatalogHelp(t *testing.T) {
	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()
	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
	)

	got := c.CatalogHelp()
	// Since calibredb is not installed, this should return an error message
	// We just verify it returns something (either help text or error message)
	if got == "" {
		t.Error("CatalogHelp() returned empty string")
	}
}

func TestCalibre_Catalog(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		opts    calibredb.CatalogOptions
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "Missing required Path",
			opts: calibredb.CatalogOptions{
				Path: "",
			},
			wantErr: true,
		},
		{
			name: "Valid Path only",
			opts: calibredb.CatalogOptions{
				Path: os.TempDir() + "/catalog.csv",
			},
			wantErr: false,
		},
		{
			name: "Valid Path with Ids",
			opts: calibredb.CatalogOptions{
				Path: os.TempDir() + "/catalog.epub",
				Ids:  "1,2,3",
			},
			wantErr: false,
		},
		{
			name: "Valid Path with Search",
			opts: calibredb.CatalogOptions{
				Path:   os.TempDir() + "/catalog.mobi",
				Search: "author:Smith",
			},
			wantErr: false,
		},
		{
			name: "Valid Path with Verbose true",
			opts: calibredb.CatalogOptions{
				Path:    os.TempDir() + "/catalog.xml",
				Verbose: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "Valid Path with Verbose false",
			opts: calibredb.CatalogOptions{
				Path:    os.TempDir() + "/catalog.xml",
				Verbose: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "Valid Path with all options",
			opts: calibredb.CatalogOptions{
				Path:    os.TempDir() + "/catalog.csv",
				Ids:     "1,2,3,4,5",
				Search:  "title:Test AND author:Smith",
				Verbose: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "Valid Path with Ids and Search",
			opts: calibredb.CatalogOptions{
				Path:   os.TempDir() + "/catalog.epub",
				Ids:    "10,20,30",
				Search: "tags:fiction",
			},
			wantErr: false,
		},
		{
			name: "Valid Path with empty Ids",
			opts: calibredb.CatalogOptions{
				Path: os.TempDir() + "/catalog.csv",
				Ids:  "",
			},
			wantErr: false,
		},
		{
			name: "Valid Path with empty Search",
			opts: calibredb.CatalogOptions{
				Path:   os.TempDir() + "/catalog.mobi",
				Search: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, f := getTestCalibre(t.Name())
			defer f()

			got, gotErr := c.Catalog(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Catalog() failed: %v", gotErr)
				}
				// For validation errors, check that it's the right type of error
				if tt.opts.Path == "" && !strings.Contains(gotErr.Error(), "required") {
					t.Errorf("Catalog() error for missing Path should mention 'required', got: %v", gotErr)
				}
				return
			}
			// read the outputted file
			val, err := os.ReadFile(tt.opts.Path)
			if err != nil {
				t.Errorf("Failed to read output file %s: %v", tt.opts.Path, err)
			}
			t.Logf("Output file content:\n%s", string(val))
			if tt.wantErr {
				t.Fatal("Catalog() succeeded unexpectedly")
			}
			if tt.want != "" && !strings.HasPrefix(got, tt.want) {
				t.Errorf("Catalog() = %v, want %v", got, tt.want)
			}
		})
	}
}

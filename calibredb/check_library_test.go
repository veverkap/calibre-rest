package calibredb_test

import (
	"os"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_CheckLibraryHelp(t *testing.T) {
	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()

	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
	)

	help := c.CheckLibraryHelp()
	if help == "" {
		t.Error("CheckLibraryHelp() returned empty string")
	}
}

func TestCalibre_CheckLibrary(t *testing.T) {
	t.Skip("TODO: Make this test work")
	tests := []struct {
		name    string // description of this test case
		opts    calibredb.CheckLibraryOptions
		args    []string
		wantErr bool
	}{
		{
			name:    "Empty options - should succeed",
			opts:    calibredb.CheckLibraryOptions{},
			wantErr: false,
		},
		{
			name: "CSV option enabled",
			opts: calibredb.CheckLibraryOptions{
				Csv: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "CSV option disabled",
			opts: calibredb.CheckLibraryOptions{
				Csv: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "IgnoreExtensions with single extension",
			opts: calibredb.CheckLibraryOptions{
				IgnoreExtensions: "txt",
			},
			wantErr: false,
		},
		{
			name: "IgnoreExtensions with multiple extensions",
			opts: calibredb.CheckLibraryOptions{
				IgnoreExtensions: "txt,pdf,doc",
			},
			wantErr: false,
		},
		{
			name: "IgnoreNames with single name",
			opts: calibredb.CheckLibraryOptions{
				IgnoreNames: "temp.txt",
			},
			wantErr: false,
		},
		{
			name: "IgnoreNames with multiple names",
			opts: calibredb.CheckLibraryOptions{
				IgnoreNames: "temp.txt,backup.db,cache.dat",
			},
			wantErr: false,
		},
		{
			name: "Report with single report type",
			opts: calibredb.CheckLibraryOptions{
				Report: "invalid_titles",
			},
			wantErr: false,
		},
		{
			name: "Report with multiple report types",
			opts: calibredb.CheckLibraryOptions{
				Report: "invalid_titles,extra_titles,missing_formats",
			},
			wantErr: false,
		},
		{
			name: "VacuumFtsDb enabled",
			opts: calibredb.CheckLibraryOptions{
				VacuumFtsDb: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "VacuumFtsDb disabled",
			opts: calibredb.CheckLibraryOptions{
				VacuumFtsDb: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "All options combined",
			opts: calibredb.CheckLibraryOptions{
				Csv:              func(b bool) *bool { return &b }(true),
				IgnoreExtensions: "txt,pdf",
				IgnoreNames:      "temp.txt,backup.db",
				Report:           "invalid_titles,extra_titles",
				VacuumFtsDb:      func(b bool) *bool { return &b }(true),
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

			_, gotErr := c.CheckLibrary(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CheckLibrary() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CheckLibrary() succeeded unexpectedly")
			}
		})
	}
}

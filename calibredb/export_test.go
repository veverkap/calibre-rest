package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_ExportHelp(t *testing.T) {
	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()
	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
	)

	got := c.ExportHelp()
	if got == "" {
		t.Error("ExportHelp() returned empty string")
	}
}

func TestCalibre_Export(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		opts    calibredb.ExportOptions
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "Missing required Ids field",
			opts: calibredb.ExportOptions{
				Ids: []string{},
			},
			wantErr: true,
		},
		{
			name: "Valid with single ID",
			opts: calibredb.ExportOptions{
				Ids: []string{"1"},
			},
			wantErr: false,
		},
		{
			name: "Valid with multiple IDs",
			opts: calibredb.ExportOptions{
				Ids: []string{"1", "2", "3"},
			},
			wantErr: false,
		},
		{
			name: "Valid with comma-separated IDs",
			opts: calibredb.ExportOptions{
				Ids: []string{"1,2,3"},
			},
			wantErr: false,
		},
		{
			name: "With All option true",
			opts: calibredb.ExportOptions{
				Ids: []string{"1"},
				All: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With All option false",
			opts: calibredb.ExportOptions{
				Ids: []string{"1"},
				All: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "With Progress option true",
			opts: calibredb.ExportOptions{
				Ids:      []string{"1"},
				Progress: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With Progress option false",
			opts: calibredb.ExportOptions{
				Ids:      []string{"1"},
				Progress: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "With SingleDir option true",
			opts: calibredb.ExportOptions{
				Ids:       []string{"1"},
				SingleDir: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With SingleDir option false",
			opts: calibredb.ExportOptions{
				Ids:       []string{"1"},
				SingleDir: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "With ToDir option",
			opts: calibredb.ExportOptions{
				Ids:   []string{"1"},
				ToDir: "/tmp/export",
			},
			wantErr: false,
		},
		{
			name: "With empty ToDir",
			opts: calibredb.ExportOptions{
				Ids:   []string{"1"},
				ToDir: "",
			},
			wantErr: false,
		},
		{
			name: "With ToDir as current directory",
			opts: calibredb.ExportOptions{
				Ids:   []string{"1"},
				ToDir: ".",
			},
			wantErr: false,
		},
		{
			name: "With ToDir as relative path",
			opts: calibredb.ExportOptions{
				Ids:   []string{"1"},
				ToDir: "./exports",
			},
			wantErr: false,
		},
		{
			name: "With ToDir as absolute path",
			opts: calibredb.ExportOptions{
				Ids:   []string{"1"},
				ToDir: "/tmp/calibre/exports",
			},
			wantErr: false,
		},
		{
			name: "With All and Progress options",
			opts: calibredb.ExportOptions{
				Ids:      []string{"1"},
				All:      func(b bool) *bool { return &b }(true),
				Progress: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With All and SingleDir options",
			opts: calibredb.ExportOptions{
				Ids:       []string{"1"},
				All:       func(b bool) *bool { return &b }(true),
				SingleDir: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With Progress and SingleDir options",
			opts: calibredb.ExportOptions{
				Ids:       []string{"1"},
				Progress:  func(b bool) *bool { return &b }(true),
				SingleDir: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With SingleDir and ToDir options",
			opts: calibredb.ExportOptions{
				Ids:       []string{"1"},
				SingleDir: func(b bool) *bool { return &b }(true),
				ToDir:     "/tmp/single-export",
			},
			wantErr: false,
		},
		{
			name: "With all options combined",
			opts: calibredb.ExportOptions{
				Ids:       []string{"1", "2", "3"},
				All:       func(b bool) *bool { return &b }(true),
				Progress:  func(b bool) *bool { return &b }(true),
				SingleDir: func(b bool) *bool { return &b }(true),
				ToDir:     "/tmp/all-exports",
			},
			wantErr: false,
		},
		{
			name: "With many IDs",
			opts: calibredb.ExportOptions{
				Ids: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			},
			wantErr: false,
		},
		{
			name: "With comma-separated and individual IDs",
			opts: calibredb.ExportOptions{
				Ids: []string{"1,2,3", "4", "5,6"},
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

			got, gotErr := c.Export(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					// Check if the error is just because calibredb is not installed
					if strings.Contains(gotErr.Error(), "no such file or directory") {
						t.Skip("Skipping test: calibredb not found")
					}
					t.Errorf("Export() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Export() succeeded unexpectedly")
			}
			if tt.want != "" && !strings.HasPrefix(got, tt.want) {
				t.Errorf("Export() = %v, want %v", got, tt.want)
			}
		})
	}
}

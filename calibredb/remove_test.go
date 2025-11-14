package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_RemoveHelp(t *testing.T) {
	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()
	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
	)

	got := c.RemoveHelp()
	// Since calibredb is not installed, this should return an error message
	// We just verify it returns something (either help text or error message)
	if got == "" {
		t.Error("RemoveHelp() returned empty string")
	}
}

func TestCalibre_Remove(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		opts    calibredb.RemoveOptions
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "Missing required Ids field",
			opts: calibredb.RemoveOptions{
				Ids: []string{},
			},
			wantErr: true,
		},
		{
			name: "Valid with single ID",
			opts: calibredb.RemoveOptions{
				Ids: []string{"1"},
			},
			wantErr: false,
		},
		{
			name: "Valid with multiple IDs",
			opts: calibredb.RemoveOptions{
				Ids: []string{"1", "2", "3"},
			},
			wantErr: false,
		},
		{
			name: "Valid with ID range",
			opts: calibredb.RemoveOptions{
				Ids: []string{"1-10"},
			},
			wantErr: false,
		},
		{
			name: "Valid with mixed IDs and ranges",
			opts: calibredb.RemoveOptions{
				Ids: []string{"1", "5-10", "15"},
			},
			wantErr: false,
		},
		{
			name: "Valid with comma separated IDs",
			opts: calibredb.RemoveOptions{
				Ids: []string{"23,34,57-85"},
			},
			wantErr: false,
		},
		{
			name: "Valid with Permanent true",
			opts: calibredb.RemoveOptions{
				Ids:       []string{"1"},
				Permanent: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "Valid with Permanent false",
			opts: calibredb.RemoveOptions{
				Ids:       []string{"1"},
				Permanent: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "Valid with Permanent nil (default)",
			opts: calibredb.RemoveOptions{
				Ids:       []string{"1"},
				Permanent: nil,
			},
			wantErr: false,
		},
		{
			name: "Valid with multiple IDs and Permanent true",
			opts: calibredb.RemoveOptions{
				Ids:       []string{"1", "2", "3", "4", "5"},
				Permanent: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "Valid with range and Permanent true",
			opts: calibredb.RemoveOptions{
				Ids:       []string{"10-20"},
				Permanent: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "Valid with mixed format and Permanent true",
			opts: calibredb.RemoveOptions{
				Ids:       []string{"1,2,3", "10-15", "20"},
				Permanent: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "Valid with all options combined",
			opts: calibredb.RemoveOptions{
				Ids:       []string{"1", "5", "10-20", "25,26,27"},
				Permanent: func(b bool) *bool { return &b }(true),
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

			got, gotErr := c.Remove(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					// Check if the error is just because calibredb is not installed
					if strings.Contains(gotErr.Error(), "no such file or directory") {
						t.Skip("Skipping test: calibredb not found")
					}
					t.Errorf("Remove() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Remove() succeeded unexpectedly")
			}
			if tt.want != "" && !strings.HasPrefix(got, tt.want) {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

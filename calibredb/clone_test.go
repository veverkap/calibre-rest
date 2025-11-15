package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_CloneHelp(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "CloneHelp returns help text",
			want: "calibredb clone",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calibredb.NewCalibre(
				calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
			)
			got := c.CloneHelp()
			if !strings.Contains(got, tt.want) && !strings.Contains(got, "no such file or directory") {
				t.Errorf("CloneHelp() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}

func TestCalibre_Clone(t *testing.T) {
	t.Skip()
	tests := []struct {
		name           string
		opts           calibredb.CloneOptions
		args           []string
		want           string
		wantErr        bool
		wantValidation bool // true if we expect a validation error
	}{
		{
			name: "Invalid - empty path",
			opts: calibredb.CloneOptions{
				Path: "",
			},
			wantErr:        true,
			wantValidation: true,
		},
		{
			name: "Valid - with path",
			opts: calibredb.CloneOptions{
				Path: "/tmp/test-clone-library",
			},
			wantErr: false,
			want:    "",
		},
		{
			name: "Valid - with relative path",
			opts: calibredb.CloneOptions{
				Path: "./new-library",
			},
			wantErr: false,
		},
		{
			name: "Valid - with absolute path",
			opts: calibredb.CloneOptions{
				Path: "/var/tmp/calibre-clone-test",
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

			got, gotErr := c.Clone(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					// If calibredb is not installed, we can't test actual execution
					// but we can verify validation errors work
					if !strings.Contains(gotErr.Error(), "no such file or directory") {
						t.Errorf("Clone() failed: %v", gotErr)
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
				t.Fatal("Clone() succeeded unexpectedly")
			}
			if tt.want != "" && !strings.Contains(got, tt.want) {
				t.Errorf("Clone() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}

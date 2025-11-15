package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_FtsIndexHelp(t *testing.T) {
	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()
	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
	)

	got := c.FtsIndexHelp()
	if got == "" {
		t.Error("FtsIndexHelp() returned empty string")
	}
}

func TestCalibre_FtsIndex(t *testing.T) {
	tests := []struct {
		name    string
		opts    calibredb.FtsIndexOptions
		args    []string
		wantErr bool
	}{
		{
			name: "Missing required EnableDisableStatusReindex field",
			opts: calibredb.FtsIndexOptions{
				EnableDisableStatusReindex: "",
			},
			wantErr: true,
		},
		{
			name: "Valid with enable",
			opts: calibredb.FtsIndexOptions{
				EnableDisableStatusReindex: "enable",
			},
			wantErr: false,
		},
		{
			name: "Valid with disable",
			opts: calibredb.FtsIndexOptions{
				EnableDisableStatusReindex: "disable",
			},
			wantErr: false,
		},
		{
			name: "Valid with status",
			opts: calibredb.FtsIndexOptions{
				EnableDisableStatusReindex: "status",
			},
			wantErr: false,
		},
		{
			name: "Valid with reindex",
			opts: calibredb.FtsIndexOptions{
				EnableDisableStatusReindex: "reindex",
			},
			wantErr: false,
		},
		{
			name: "With IndexingSpeed fast",
			opts: calibredb.FtsIndexOptions{
				EnableDisableStatusReindex: "reindex",
				IndexingSpeed:              calibredb.Fast,
			},
			wantErr: false,
		},
		{
			name: "With IndexingSpeed slow",
			opts: calibredb.FtsIndexOptions{
				EnableDisableStatusReindex: "enable",
				IndexingSpeed:              calibredb.Slow,
			},
			wantErr: false,
		},
		{
			name: "With WaitForCompletion true",
			opts: calibredb.FtsIndexOptions{
				EnableDisableStatusReindex: "reindex",
				WaitForCompletion:          func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With WaitForCompletion false",
			opts: calibredb.FtsIndexOptions{
				EnableDisableStatusReindex: "enable",
				WaitForCompletion:          func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "With all options",
			opts: calibredb.FtsIndexOptions{
				EnableDisableStatusReindex: "reindex",
				IndexingSpeed:              calibredb.Fast,
				WaitForCompletion:          func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "With IndexingSpeed and WaitForCompletion false",
			opts: calibredb.FtsIndexOptions{
				EnableDisableStatusReindex: "enable",
				IndexingSpeed:              calibredb.Slow,
				WaitForCompletion:          func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "Empty IndexingSpeed should be omitted",
			opts: calibredb.FtsIndexOptions{
				EnableDisableStatusReindex: "status",
				IndexingSpeed:              "",
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

			got, gotErr := c.FtsIndex(tt.opts, tt.args...)
			if gotErr != nil {
				// If calibredb is not installed, we expect an error
				// The test validates that the function can be called properly
				if !strings.Contains(gotErr.Error(), "no such file or directory") &&
					!strings.Contains(gotErr.Error(), "executable file not found") {
					if !tt.wantErr {
						t.Errorf("FtsIndex() failed with unexpected error: %v", gotErr)
					}
				} else {
					// It's a "calibredb not found" error
					if tt.wantErr && strings.Contains(gotErr.Error(), "required") {
						// This is a validation error, which is expected
						return
					}
					// For missing calibredb, we skip further checks if we don't expect validation errors
					if !tt.wantErr {
						return
					}
				}
				// Check for validation errors
				if tt.opts.EnableDisableStatusReindex == "" && !strings.Contains(gotErr.Error(), "required") {
					t.Errorf("FtsIndex() error for missing EnableDisableStatusReindex should mention 'required', got: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("FtsIndex() succeeded unexpectedly")
			}
			// If we somehow have calibredb installed and it succeeds,
			// verify we got some output
			if got == "" && gotErr == nil {
				t.Error("FtsIndex() returned empty string without error")
			}
		})
	}
}

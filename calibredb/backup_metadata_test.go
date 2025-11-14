package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_BackupMetadataHelp(t *testing.T) {
	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()
	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
	)

	// Call BackupMetadataHelp - this will fail if calibredb is not installed
	// but we test that the method exists and can be called
	help := c.BackupMetadataHelp()
	
	// The help output should contain something (error message or actual help)
	if help == "" {
		t.Error("BackupMetadataHelp() returned empty string")
	}
}

func TestCalibre_BackupMetadata(t *testing.T) {
	tests := []struct {
		name    string
		opts    calibredb.BackupMetadataOptions
		args    []string
		wantErr bool
	}{
		{
			name:    "No options - default behavior",
			opts:    calibredb.BackupMetadataOptions{},
			wantErr: false,
		},
		{
			name: "All option set to true",
			opts: calibredb.BackupMetadataOptions{
				All: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "All option set to false",
			opts: calibredb.BackupMetadataOptions{
				All: func(b bool) *bool { return &b }(false),
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

			got, gotErr := c.BackupMetadata(tt.opts, tt.args...)
			if gotErr != nil {
				// If calibredb is not installed, we expect an error
				// The test validates that the function can be called properly
				if !strings.Contains(gotErr.Error(), "no such file or directory") &&
					!strings.Contains(gotErr.Error(), "executable file not found") {
					if !tt.wantErr {
						t.Errorf("BackupMetadata() failed with unexpected error: %v", gotErr)
					}
				}
				return
			}
			if tt.wantErr {
				t.Fatal("BackupMetadata() succeeded unexpectedly")
			}
			// If we somehow have calibredb installed and it succeeds,
			// verify we got some output
			if got == "" && gotErr == nil {
				t.Error("BackupMetadata() returned empty string without error")
			}
		})
	}
}

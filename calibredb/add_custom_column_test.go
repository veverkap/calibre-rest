package calibredb_test

import (
	"os"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_AddCustomColumn(t *testing.T) {
	tempDir := os.TempDir()
	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
	)
	tests := []struct {
		name    string // description of this test case
		opts    calibredb.AddCustomColumnOptions
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "Invalid label",
			opts: calibredb.AddCustomColumnOptions{
				Label:    "",
				Name:     "My Column",
				Datatype: "text",
			},
			wantErr: true,
		},
		{
			name: "Invalid name",
			opts: calibredb.AddCustomColumnOptions{
				Label:    "my_column",
				Name:     "",
				Datatype: "text",
			},
			wantErr: true,
		},
		{
			name: "Empty datatype",
			opts: calibredb.AddCustomColumnOptions{
				Label:    "my_column",
				Name:     "My Column",
				Datatype: "",
			},
			wantErr: true,
		},
		{
			name: "Invalid datatype",
			opts: calibredb.AddCustomColumnOptions{
				Label:    "my_column",
				Name:     "My Column",
				Datatype: "invalid_type",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := c.AddCustomColumn(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("AddCustomColumn() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("AddCustomColumn() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("AddCustomColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}

package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_AddCustomColumn(t *testing.T) {
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
		{
			name: "Valid custom column",
			opts: calibredb.AddCustomColumnOptions{
				Label:    "validito",
				Name:     "my column",
				Datatype: "text",
			},
			wantErr: false,
			want:    "Custom column created with id",
		},
		{
			name: "Valid custom column with options",
			opts: calibredb.AddCustomColumnOptions{
				Label:      "valid_with_options",
				Name:       "specific_column",
				Datatype:   "text",
				Display:    `{"is_names": true, "use_decorations": false}`,
				IsMultiple: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
			want:    "Custom column created with id",
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
			if !strings.HasPrefix(got, tt.want) {
				t.Errorf("AddCustomColumn() = %v, want %v", got, tt.want)
			}
			custom, err := c.CustomColumns(
				calibredb.CustomColumnsOptions{
					Details: func(b bool) *bool { return &b }(true),
				},
			)
			if err != nil {
				t.Fatalf("CustomColumns() failed: %v", err)
			}
			if !strings.Contains(custom, tt.opts.Label) || !strings.Contains(custom, tt.opts.Name) {
				t.Errorf("CustomColumns() = %v, want it to contain label %v and name %v", custom, tt.opts.Label, tt.opts.Name)
			}
		})
	}
}

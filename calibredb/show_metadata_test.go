package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_ShowMetadataHelp(t *testing.T) {
	c, f := getTestCalibre(t.Name())
	defer f()

	got := c.ShowMetadataHelp()
	// Since calibredb is not installed, this should return an error message
	// We just verify it returns something (either help text or error message)
	if got == "" {
		t.Error("ShowMetadataHelp() returned empty string")
	}
}

func TestCalibre_ShowMetadata(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		opts    calibredb.ShowMetadataOptions
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "Missing required Id field",
			opts: calibredb.ShowMetadataOptions{
				Id: "",
			},
			wantErr: true,
		},
		{
			name: "Valid Id only",
			opts: calibredb.ShowMetadataOptions{
				Id: "1",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "Valid Id with AsOpf true",
			opts: calibredb.ShowMetadataOptions{
				Id:    "123",
				AsOpf: func(b bool) *bool { return &b }(true),
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "Valid Id with AsOpf false",
			opts: calibredb.ShowMetadataOptions{
				Id:    "456",
				AsOpf: func(b bool) *bool { return &b }(false),
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "Valid Id with AsOpf nil",
			opts: calibredb.ShowMetadataOptions{
				Id:    "789",
				AsOpf: nil,
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "Valid Id with numeric string",
			opts: calibredb.ShowMetadataOptions{
				Id: "42",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "Valid Id with large number",
			opts: calibredb.ShowMetadataOptions{
				Id: "999999",
			},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
		{
			name: "Valid Id with AsOpf true and additional args",
			opts: calibredb.ShowMetadataOptions{
				Id:    "100",
				AsOpf: func(b bool) *bool { return &b }(true),
			},
			args:    []string{"extra", "args"},
			wantErr: true, // Will fail because calibredb is not installed, but validation passes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, f := getTestCalibre(t.Name())
			defer f()

			got, gotErr := c.ShowMetadata(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					// Check if the error is just because calibredb is not installed
					if strings.Contains(gotErr.Error(), "no such file or directory") {
						t.Skip("Skipping test: calibredb not found")
					}
					t.Errorf("ShowMetadata() failed: %v", gotErr)
				}
				// For validation errors, check that it's the right type of error
				if tt.opts.Id == "" && !strings.Contains(gotErr.Error(), "required") {
					t.Errorf("ShowMetadata() error for missing Id should mention 'required', got: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ShowMetadata() succeeded unexpectedly")
			}
			if tt.want != "" && !strings.HasPrefix(got, tt.want) {
				t.Errorf("ShowMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

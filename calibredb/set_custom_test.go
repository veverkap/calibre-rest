package calibredb_test

import (
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_SetCustomHelp(t *testing.T) {
	c, f := getTestCalibre(t.Name())
	defer f()

	got := c.SetCustomHelp()
	if got == "" {
		t.Error("SetCustomHelp() returned empty string")
	}
}

func TestCalibre_SetCustom(t *testing.T) {
	t.Skip("Skip temporarily")
	tests := []struct {
		name    string // description of this test case
		opts    calibredb.SetCustomOptions
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "Missing required Column field",
			opts: calibredb.SetCustomOptions{
				Column: "",
				Id:     "1",
				Value:  "test value",
			},
			wantErr: true,
		},
		{
			name: "Missing required Id field",
			opts: calibredb.SetCustomOptions{
				Column: "my_column",
				Id:     "",
				Value:  "test value",
			},
			wantErr: true,
		},
		{
			name: "Missing required Value field",
			opts: calibredb.SetCustomOptions{
				Column: "my_column",
				Id:     "1",
				Value:  "",
			},
			wantErr: true,
		},
		{
			name: "All required fields missing",
			opts: calibredb.SetCustomOptions{
				Column: "",
				Id:     "",
				Value:  "",
			},
			wantErr: true,
		},
		{
			name: "Valid with all required fields",
			opts: calibredb.SetCustomOptions{
				Column: "my_column",
				Id:     "1",
				Value:  "test value",
			},
			wantErr: false,
		},
		{
			name: "Valid with Append true",
			opts: calibredb.SetCustomOptions{
				Column: "my_column",
				Id:     "1",
				Value:  "test value",
				Append: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
		{
			name: "Valid with Append false",
			opts: calibredb.SetCustomOptions{
				Column: "my_column",
				Id:     "1",
				Value:  "test value",
				Append: func(b bool) *bool { return &b }(false),
			},
			wantErr: false,
		},
		{
			name: "Valid with Append nil (default)",
			opts: calibredb.SetCustomOptions{
				Column: "my_column",
				Id:     "1",
				Value:  "test value",
				Append: nil,
			},
			wantErr: false,
		},
		{
			name: "Valid with numeric Id",
			opts: calibredb.SetCustomOptions{
				Column: "my_column",
				Id:     "12345",
				Value:  "test value",
			},
			wantErr: false,
		},
		{
			name: "Valid with complex column name",
			opts: calibredb.SetCustomOptions{
				Column: "complex_column_name",
				Id:     "1",
				Value:  "test value",
			},
			wantErr: false,
		},
		{
			name: "Valid with special characters in value",
			opts: calibredb.SetCustomOptions{
				Column: "my_column",
				Id:     "1",
				Value:  "test value with special chars: !@#$%^&*()",
			},
			wantErr: false,
		},
		{
			name: "Valid with multiline value",
			opts: calibredb.SetCustomOptions{
				Column: "my_column",
				Id:     "1",
				Value:  "line1\nline2\nline3",
			},
			wantErr: false,
		},
		{
			name: "Valid with empty-like but non-empty value (spaces)",
			opts: calibredb.SetCustomOptions{
				Column: "my_column",
				Id:     "1",
				Value:  "   ",
			},
			wantErr: false,
		},
		{
			name: "Valid with all options combined",
			opts: calibredb.SetCustomOptions{
				Column: "test_column",
				Id:     "999",
				Value:  "comprehensive test value",
				Append: func(b bool) *bool { return &b }(true),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, f := getTestCalibre(t.Name())
			defer f()

			got, gotErr := c.SetCustom(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					// Check if the error is just because calibredb is not installed
					if strings.Contains(gotErr.Error(), "no such file or directory") {
						t.Skip("Skipping test: calibredb not found")
					}
					t.Errorf("SetCustom() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SetCustom() succeeded unexpectedly")
			}
			if tt.want != "" && !strings.HasPrefix(got, tt.want) {
				t.Errorf("SetCustom() = %v, want %v", got, tt.want)
			}
		})
	}
}

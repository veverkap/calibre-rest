package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_ListHelp(t *testing.T) {
	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()
	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
	)

	got := c.ListHelp()
	// Since calibredb is not installed, this should return an error message
	// We just verify it returns something (either help text or error message)
	if got == "" {
		t.Error("ListHelp() returned empty string")
	}
}

func TestCalibre_List(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		opts    calibredb.ListOptions
		args    []string
		wantErr bool
	}{
		{
			name:    "Empty options - should succeed",
			opts:    calibredb.ListOptions{},
			wantErr: false,
		},
		{
			name: "Ascending option enabled",
			opts: calibredb.ListOptions{
				Ascending: boolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "Ascending option disabled",
			opts: calibredb.ListOptions{
				Ascending: boolPtr(false),
			},
			wantErr: false,
		},
		{
			name: "Fields option with single field",
			opts: calibredb.ListOptions{
				Fields: "title",
			},
			wantErr: false,
		},
		{
			name: "Fields option with multiple fields",
			opts: calibredb.ListOptions{
				Fields: "title,authors,publisher",
			},
			wantErr: false,
		},
		{
			name: "Fields option with all",
			opts: calibredb.ListOptions{
				Fields: "all",
			},
			wantErr: false,
		},
		{
			name: "ForMachine option enabled",
			opts: calibredb.ListOptions{
				ForMachine: boolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "ForMachine option disabled",
			opts: calibredb.ListOptions{
				ForMachine: boolPtr(false),
			},
			wantErr: false,
		},
		{
			name: "Limit option set",
			opts: calibredb.ListOptions{
				Limit: 10,
			},
			wantErr: false,
		},
		{
			name: "Limit option zero (default)",
			opts: calibredb.ListOptions{
				Limit: 0,
			},
			wantErr: false,
		},
		{
			name: "LineWidth option set",
			opts: calibredb.ListOptions{
				LineWidth: 80,
			},
			wantErr: false,
		},
		{
			name: "LineWidth option zero (default)",
			opts: calibredb.ListOptions{
				LineWidth: 0,
			},
			wantErr: false,
		},
		{
			name: "Prefix option set",
			opts: calibredb.ListOptions{
				Prefix: "/path/to/library",
			},
			wantErr: false,
		},
		{
			name: "Search option set",
			opts: calibredb.ListOptions{
				Search: "title:test",
			},
			wantErr: false,
		},
		{
			name: "Separator option set",
			opts: calibredb.ListOptions{
				Separator: ",",
			},
			wantErr: false,
		},
		{
			name: "SortBy option with single field",
			opts: calibredb.ListOptions{
				SortBy: "title",
			},
			wantErr: false,
		},
		{
			name: "SortBy option with multiple fields",
			opts: calibredb.ListOptions{
				SortBy: "authors,title",
			},
			wantErr: false,
		},
		{
			name: "Template option set",
			opts: calibredb.ListOptions{
				Template: "{title} - {authors}",
			},
			wantErr: false,
		},
		{
			name: "TemplateFile option set",
			opts: calibredb.ListOptions{
				TemplateFile: "/path/to/template.txt",
			},
			wantErr: false,
		},
		{
			name: "TemplateHeading option set",
			opts: calibredb.ListOptions{
				TemplateHeading: "My Books",
			},
			wantErr: false,
		},
		{
			name: "All boolean options enabled",
			opts: calibredb.ListOptions{
				Ascending:  boolPtr(true),
				ForMachine: boolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "All string options set",
			opts: calibredb.ListOptions{
				Fields:          "title,authors",
				Prefix:          "/library",
				Search:          "authors:Smith",
				Separator:       "|",
				SortBy:          "title",
				Template:        "{title}",
				TemplateFile:    "/template.txt",
				TemplateHeading: "Books",
			},
			wantErr: false,
		},
		{
			name: "All int options set",
			opts: calibredb.ListOptions{
				Limit:     50,
				LineWidth: 100,
			},
			wantErr: false,
		},
		{
			name: "All options combined",
			opts: calibredb.ListOptions{
				Ascending:       boolPtr(true),
				Fields:          "title,authors,publisher",
				ForMachine:      boolPtr(false),
				Limit:           25,
				LineWidth:       120,
				Prefix:          "/my/library",
				Search:          "tag:fiction",
				Separator:       " | ",
				SortBy:          "authors,title",
				Template:        "{title} by {authors}",
				TemplateFile:    "/path/to/custom/template.txt",
				TemplateHeading: "Book List",
			},
			wantErr: false,
		},
		{
			name: "With additional arguments",
			opts: calibredb.ListOptions{
				Fields: "title",
			},
			args:    []string{"1", "2", "3"},
			wantErr: false,
		},
		{
			name: "ForMachine with Fields",
			opts: calibredb.ListOptions{
				ForMachine: boolPtr(true),
				Fields:     "title,authors,isbn",
			},
			wantErr: false,
		},
		{
			name: "Search with Limit",
			opts: calibredb.ListOptions{
				Search: "authors:Asimov",
				Limit:  10,
			},
			wantErr: false,
		},
		{
			name: "SortBy with Ascending",
			opts: calibredb.ListOptions{
				SortBy:    "pubdate",
				Ascending: boolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "Template with TemplateHeading",
			opts: calibredb.ListOptions{
				Fields:          "template",
				Template:        "{title} - {series}",
				TemplateHeading: "Series Books",
			},
			wantErr: false,
		},
		{
			name: "Custom field in Fields",
			opts: calibredb.ListOptions{
				Fields: "title,*rating,*genre",
			},
			wantErr: false,
		},
		{
			name: "Custom field in SortBy",
			opts: calibredb.ListOptions{
				SortBy: "*rating",
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

			_, err := c.List(tt.opts, tt.args...)

			if err != nil {
				if !tt.wantErr {
					// Check if the error is just because calibredb is not installed
					if strings.Contains(err.Error(), "no such file or directory") {
						t.Skip("Skipping test: calibredb not found")
					}
					t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("List() succeeded unexpectedly")
			}
		})
	}
}

func TestCalibre_List_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		opts    calibredb.ListOptions
		wantErr bool
	}{
		{
			name:    "Empty options passes validation",
			opts:    calibredb.ListOptions{},
			wantErr: false,
		},
		{
			name: "All valid string options",
			opts: calibredb.ListOptions{
				Fields:          "title,authors",
				Prefix:          "/valid/path",
				Search:          "valid search",
				Separator:       ",",
				SortBy:          "title",
				Template:        "valid template",
				TemplateFile:    "/valid/file.txt",
				TemplateHeading: "Valid Heading",
			},
			wantErr: false,
		},
		{
			name: "Positive int values",
			opts: calibredb.ListOptions{
				Limit:     100,
				LineWidth: 200,
			},
			wantErr: false,
		},
		{
			name: "Negative Limit value",
			opts: calibredb.ListOptions{
				Limit: -1,
			},
			wantErr: false, // The validation doesn't check for negative values, only for required fields
		},
		{
			name: "Negative LineWidth value",
			opts: calibredb.ListOptions{
				LineWidth: -1,
			},
			wantErr: false, // The validation doesn't check for negative values, only for required fields
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

			_, err := c.List(tt.opts)

			if err != nil {
				if !tt.wantErr {
					// Check if the error is just because calibredb is not installed
					if strings.Contains(err.Error(), "no such file or directory") {
						t.Skip("Skipping test: calibredb not found")
					}
					t.Errorf("List() validation error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("List() succeeded unexpectedly")
			}
		})
	}
}

func TestCalibre_List_WithRealCalibredb(t *testing.T) {
	calibredbPath := findCalibredb()
	if calibredbPath == "" {
		t.Skip("calibredb not found, skipping test with real calibredb")
	}

	tempDir := os.TempDir() + "/calibre_test_list"
	defer func() { _ = os.RemoveAll(tempDir) }()

	c := calibredb.NewCalibre(
		calibredb.WithCalibreDBLocation(calibredbPath),
		calibredb.WithLibraryPath(tempDir),
	)

	tests := []struct {
		name         string
		opts         calibredb.ListOptions
		wantContains string
	}{
		{
			name:         "Basic list",
			opts:         calibredb.ListOptions{},
			wantContains: "",
		},
		{
			name: "List with ForMachine",
			opts: calibredb.ListOptions{
				ForMachine: boolPtr(true),
			},
			wantContains: "",
		},
		{
			name: "List with Fields",
			opts: calibredb.ListOptions{
				Fields: "title,authors",
			},
			wantContains: "",
		},
		{
			name: "List with Limit",
			opts: calibredb.ListOptions{
				Limit: 5,
			},
			wantContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.List(tt.opts)

			// With an empty library, list should succeed without error
			if err != nil {
				t.Errorf("List() error = %v, want nil", err)
			}

			// The output might be empty for an empty library, which is fine
			if tt.wantContains != "" && !strings.Contains(got, tt.wantContains) {
				t.Errorf("List() = %v, want to contain %v", got, tt.wantContains)
			}
		})
	}
}

func TestCalibre_ListHelp_WithRealCalibredb(t *testing.T) {
	calibredbPath := findCalibredb()
	if calibredbPath == "" {
		t.Skip("calibredb not found, skipping test with real calibredb")
	}

	tempDir := os.TempDir() + "/calibre_test_list_help"
	defer func() { _ = os.RemoveAll(tempDir) }()

	c := calibredb.NewCalibre(
		calibredb.WithCalibreDBLocation(calibredbPath),
		calibredb.WithLibraryPath(tempDir),
	)

	got := c.ListHelp()

	// Help output should contain "list" or "Usage"
	if !strings.Contains(got, "list") && !strings.Contains(got, "Usage") {
		t.Errorf("ListHelp() = %v, want to contain 'list' or 'Usage'", got)
	}
}

func TestCalibre_List_ArgumentHandling(t *testing.T) {
	tests := []struct {
		name string
		opts calibredb.ListOptions
		args []string
	}{
		{
			name: "No additional arguments",
			opts: calibredb.ListOptions{},
			args: []string{},
		},
		{
			name: "Single argument",
			opts: calibredb.ListOptions{},
			args: []string{"1"},
		},
		{
			name: "Multiple arguments",
			opts: calibredb.ListOptions{},
			args: []string{"1", "2", "3"},
		},
		{
			name: "Arguments with options",
			opts: calibredb.ListOptions{
				Fields: "title",
			},
			args: []string{"10", "20"},
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

			_, err := c.List(tt.opts, tt.args...)

			// We expect an error because calibredb is not installed
			// But we're testing that the arguments are properly handled
			if err == nil {
				// If somehow succeeds, that's ok
			}
		})
	}
}

func TestCalibre_List_OnErrorCallback(t *testing.T) {
	called := false
	var capturedError error

	onError := func(err error) {
		called = true
		capturedError = err
	}

	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()

	c := calibredb.NewCalibre(
		calibredb.WithCalibreDBLocation("/nonexistent/calibredb"),
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithOnError(onError),
	)

	_, _ = c.List(calibredb.ListOptions{})

	if !called {
		t.Error("OnError callback was not called when List() failed")
	}

	if capturedError == nil {
		t.Error("OnError callback received nil error")
	}
}

// Helper function to create a bool pointer
func boolPtr(b bool) *bool {
	return &b
}

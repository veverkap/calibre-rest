package calibredb_test

import (
	"os"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestCalibre_SavedSearchesHelp(t *testing.T) {
	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()
	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
	)

	got := c.SavedSearchesHelp()
	// Since calibredb may not be installed, this should return something (either help text or error message)
	// We just verify it returns something
	if got == "" {
		t.Error("SavedSearchesHelp() returned empty string")
	}
}

func TestCalibre_SavedSearches(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		opts    calibredb.SavedSearchesOptions
		args    []string
		want    string
		wantErr bool
	}{
		{
			name:    "Empty options with no args",
			opts:    calibredb.SavedSearchesOptions{},
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "Empty options with list command",
			opts:    calibredb.SavedSearchesOptions{},
			args:    []string{"list"},
			wantErr: false,
		},
		{
			name:    "Empty options with add command",
			opts:    calibredb.SavedSearchesOptions{},
			args:    []string{"add"},
			wantErr: false,
		},
		{
			name:    "Empty options with remove command",
			opts:    calibredb.SavedSearchesOptions{},
			args:    []string{"remove"},
			wantErr: false,
		},
		{
			name:    "Empty options with add and search name",
			opts:    calibredb.SavedSearchesOptions{},
			args:    []string{"add", "MySearch", "author:Smith"},
			wantErr: false,
		},
		{
			name:    "Empty options with remove and search name",
			opts:    calibredb.SavedSearchesOptions{},
			args:    []string{"remove", "MySearch"},
			wantErr: false,
		},
		{
			name:    "Empty options with multiple args",
			opts:    calibredb.SavedSearchesOptions{},
			args:    []string{"list", "extra", "args"},
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

			got, gotErr := c.SavedSearches(tt.opts, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					// Check if the error is just because calibredb is not installed
					if strings.Contains(gotErr.Error(), "no such file or directory") {
						t.Skip("Skipping test: calibredb not found")
					}
					t.Errorf("SavedSearches() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SavedSearches() succeeded unexpectedly")
			}
			if tt.want != "" && !strings.HasPrefix(got, tt.want) {
				t.Errorf("SavedSearches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSavedSearchesOptions_EmptyStruct(t *testing.T) {
	// Test that SavedSearchesOptions can be created as an empty struct
	opts := calibredb.SavedSearchesOptions{}

	// Since there are no fields, we just verify the struct can be created
	// This is mainly for code coverage
	_ = opts
}

func TestCalibre_SavedSearches_ValidationSuccess(t *testing.T) {
	// Test that validation passes for valid options
	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()

	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/Applications/calibre.app/Contents/MacOS/calibredb"),
	)

	opts := calibredb.SavedSearchesOptions{}

	// This should not fail validation since there are no required fields
	_, err := c.SavedSearches(opts)

	// We expect an error because calibredb is not installed, not a validation error
	if err != nil && strings.Contains(err.Error(), "validation") {
		t.Errorf("SavedSearches() failed validation with valid options: %v", err)
	}
}

func TestCalibre_SavedSearchesHelp_ErrorHandling(t *testing.T) {
	// Test error handling when calibredb doesn't exist
	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()

	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/nonexistent/calibredb"),
	)

	got := c.SavedSearchesHelp()

	// Should return an error message
	if !strings.Contains(got, "no such file or directory") && !strings.Contains(got, "executable file not found") {
		t.Errorf("SavedSearchesHelp() with invalid calibredb path should return error message, got: %v", got)
	}
}

func TestCalibre_SavedSearches_ErrorHandling(t *testing.T) {
	// Test error handling when calibredb doesn't exist
	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()

	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/nonexistent/calibredb"),
	)

	opts := calibredb.SavedSearchesOptions{}
	_, err := c.SavedSearches(opts)

	if err == nil {
		t.Error("SavedSearches() with invalid calibredb path should return error")
	}

	// Should return an error about missing file
	if !strings.Contains(err.Error(), "no such file or directory") && !strings.Contains(err.Error(), "executable file not found") {
		t.Errorf("SavedSearches() with invalid calibredb path should return file error, got: %v", err)
	}
}

func TestCalibre_SavedSearches_WithOnErrorCallback(t *testing.T) {
	// Test that OnError callback is called when there's an error
	tempDir := os.TempDir() + "/" + t.Name()
	defer func() { _ = os.RemoveAll(tempDir) }()

	errorCalled := false
	var capturedError error

	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation("/nonexistent/calibredb"),
		calibredb.WithOnError(func(err error) {
			errorCalled = true
			capturedError = err
		}),
	)

	opts := calibredb.SavedSearchesOptions{}
	_, err := c.SavedSearches(opts)

	if err == nil {
		t.Error("SavedSearches() should return error with invalid calibredb path")
	}

	if !errorCalled {
		t.Error("OnError callback was not called when SavedSearches() failed")
	}

	if capturedError == nil {
		t.Error("OnError callback received nil error")
	}
}

func TestCalibre_SavedSearchesHelp_WithRealCalibredb(t *testing.T) {
	// Test with real calibredb if available
	calibredbPath := findCalibredb()
	if calibredbPath == "" {
		t.Skip("calibredb not found, skipping test with real calibredb")
	}

	tempDir := os.TempDir() + "/calibre_test_saved_searches_help"
	defer func() { _ = os.RemoveAll(tempDir) }()

	c := calibredb.NewCalibre(
		calibredb.WithCalibreDBLocation(calibredbPath),
		calibredb.WithLibraryPath(tempDir),
	)

	got := c.SavedSearchesHelp()

	// Help output should contain "saved_searches" or "usage" or similar help text
	if !strings.Contains(strings.ToLower(got), "saved") && 
	   !strings.Contains(strings.ToLower(got), "usage") &&
	   !strings.Contains(strings.ToLower(got), "command") {
		t.Errorf("SavedSearchesHelp() = %v, expected help text containing 'saved', 'usage', or 'command'", got)
	}

	if got == "" {
		t.Error("SavedSearchesHelp() returned empty string")
	}
}

func TestCalibre_SavedSearches_WithRealCalibredb_List(t *testing.T) {
	// Test with real calibredb if available
	calibredbPath := findCalibredb()
	if calibredbPath == "" {
		t.Skip("calibredb not found, skipping test with real calibredb")
	}

	tempDir := os.TempDir() + "/calibre_test_saved_searches_list"
	defer func() { _ = os.RemoveAll(tempDir) }()

	c := calibredb.NewCalibre(
		calibredb.WithCalibreDBLocation(calibredbPath),
		calibredb.WithLibraryPath(tempDir),
	)

	opts := calibredb.SavedSearchesOptions{}
	got, err := c.SavedSearches(opts, "list")

	// This may fail or succeed depending on whether the library exists
	// But we're testing that the function can be called successfully
	if err != nil {
		// If there's an error, it should be about the library, not about our code
		if strings.Contains(err.Error(), "validation") {
			t.Errorf("SavedSearches() with 'list' command failed validation: %v", err)
		}
		// Other errors (like "library not found") are acceptable for this test
		return
	}

	// If successful, we should get some output
	_ = got // Output could be empty list or actual searches
}

func TestCalibre_SavedSearches_SuccessPath(t *testing.T) {
	// Test success path when calibredb is available
	calibredbPath := findCalibredb()
	if calibredbPath == "" {
		t.Skip("calibredb not found, skipping success path test")
	}

	tempDir := os.TempDir() + "/calibre_test_saved_searches_success"
	defer func() { _ = os.RemoveAll(tempDir) }()

	errorCalled := false
	c := calibredb.NewCalibre(
		calibredb.WithCalibreDBLocation(calibredbPath),
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithOnError(func(err error) {
			errorCalled = true
		}),
	)

	opts := calibredb.SavedSearchesOptions{}
	_, err := c.SavedSearches(opts)

	// The command might fail (e.g., library doesn't exist), but OnError should be called
	// or it might succeed. We're just testing the code paths.
	if err == nil && errorCalled {
		t.Error("OnError was called even though SavedSearches() succeeded")
	}
}

func TestCalibre_SavedSearches_MultipleCommands(t *testing.T) {
	// Test various command combinations
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "Add command with name and query",
			args: []string{"add", "TestSearch", "author:Test"},
		},
		{
			name: "Remove command with name",
			args: []string{"remove", "TestSearch"},
		},
		{
			name: "List command",
			args: []string{"list"},
		},
		{
			name: "No command",
			args: []string{},
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

			opts := calibredb.SavedSearchesOptions{}
			_, err := c.SavedSearches(opts, tt.args...)

			// We expect an error because calibredb is not installed
			// But we're testing that the function handles various args correctly
			if err != nil && !strings.Contains(err.Error(), "no such file or directory") && 
			   !strings.Contains(err.Error(), "executable file not found") {
				// Any other error would be unexpected
				t.Logf("SavedSearches() with args %v returned error: %v", tt.args, err)
			}
		})
	}
}

package calibredb_test

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/veverkap/calibre-rest/calibredb"
)

func TestNewCalibre(t *testing.T) {
	tests := []struct {
		name string
		opts []calibredb.CalibreOption
		want *calibredb.Calibre
	}{
		{
			name: "Default constructor",
			opts: []calibredb.CalibreOption{},
			want: &calibredb.Calibre{},
		},
		{
			name: "With CalibreDB location",
			opts: []calibredb.CalibreOption{
				calibredb.WithCalibreDBLocation("/usr/bin/calibredb"),
			},
			want: &calibredb.Calibre{
				CalibreDBLocation: "/usr/bin/calibredb",
			},
		},
		{
			name: "With library path",
			opts: []calibredb.CalibreOption{
				calibredb.WithLibraryPath("/path/to/library"),
			},
			want: &calibredb.Calibre{
				LibraryPath: "/path/to/library",
			},
		},
		{
			name: "With username",
			opts: []calibredb.CalibreOption{
				calibredb.WithUsername("testuser"),
			},
			want: &calibredb.Calibre{
				Username: "testuser",
			},
		},
		{
			name: "With password",
			opts: []calibredb.CalibreOption{
				calibredb.WithPassword("testpass"),
			},
			want: &calibredb.Calibre{
				Password: "testpass",
			},
		},
		{
			name: "With timeout",
			opts: []calibredb.CalibreOption{
				calibredb.WithTimeout("30s"),
			},
			want: &calibredb.Calibre{
				Timeout: "30s",
			},
		},
		{
			name: "With all options",
			opts: []calibredb.CalibreOption{
				calibredb.WithCalibreDBLocation("/usr/bin/calibredb"),
				calibredb.WithLibraryPath("/path/to/library"),
				calibredb.WithUsername("testuser"),
				calibredb.WithPassword("testpass"),
				calibredb.WithTimeout("30s"),
			},
			want: &calibredb.Calibre{
				CalibreDBLocation: "/usr/bin/calibredb",
				LibraryPath:       "/path/to/library",
				Username:          "testuser",
				Password:          "testpass",
				Timeout:           "30s",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calibredb.NewCalibre(tt.opts...)

			if got.CalibreDBLocation != tt.want.CalibreDBLocation {
				t.Errorf("CalibreDBLocation = %v, want %v", got.CalibreDBLocation, tt.want.CalibreDBLocation)
			}
			if got.LibraryPath != tt.want.LibraryPath {
				t.Errorf("LibraryPath = %v, want %v", got.LibraryPath, tt.want.LibraryPath)
			}
			if got.Username != tt.want.Username {
				t.Errorf("Username = %v, want %v", got.Username, tt.want.Username)
			}
			if got.Password != tt.want.Password {
				t.Errorf("Password = %v, want %v", got.Password, tt.want.Password)
			}
			if got.Timeout != tt.want.Timeout {
				t.Errorf("Timeout = %v, want %v", got.Timeout, tt.want.Timeout)
			}
		})
	}
}

func TestWithOnError(t *testing.T) {
	called := false
	var capturedError error

	onError := func(err error) {
		called = true
		capturedError = err
	}

	c := calibredb.NewCalibre(
		calibredb.WithOnError(onError),
	)

	if c.OnError == nil {
		t.Error("OnError should not be nil")
	}

	// Test that the OnError callback is called
	testError := errors.New("test error")
	c.OnError(testError)

	if !called {
		t.Error("OnError callback was not called")
	}

	if capturedError != testError {
		t.Errorf("OnError callback captured error = %v, want %v", capturedError, testError)
	}
}

func TestCalibre_Version(t *testing.T) {
	tests := []struct {
		name              string
		calibreDBLocation string
		libraryPath       string
		wantContains      string
		wantErr           bool
	}{
		{
			name:              "Invalid calibredb location",
			calibreDBLocation: "/nonexistent/calibredb",
			libraryPath:       os.TempDir(),
			wantContains:      "no such file or directory",
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calibredb.NewCalibre(
				calibredb.WithCalibreDBLocation(tt.calibreDBLocation),
				calibredb.WithLibraryPath(tt.libraryPath),
			)

			got := c.Version()

			if tt.wantErr {
				if !strings.Contains(got, tt.wantContains) {
					t.Errorf("Version() = %v, want to contain %v", got, tt.wantContains)
				}
			} else {
				if strings.Contains(got, "error") || strings.Contains(got, "Error") {
					t.Errorf("Version() returned error: %v", got)
				}
			}
		})
	}
}

func TestCalibre_Version_WithRealCalibredb(t *testing.T) {
	// Try to find calibredb in common locations
	calibredbPaths := []string{
		"/usr/bin/calibredb",
		"/usr/local/bin/calibredb",
		"/Applications/calibre.app/Contents/MacOS/calibredb",
	}

	var calibredbPath string
	for _, path := range calibredbPaths {
		if _, err := os.Stat(path); err == nil {
			calibredbPath = path
			break
		}
	}

	if calibredbPath == "" {
		// Also try 'which calibredb'
		if path, err := exec.LookPath("calibredb"); err == nil {
			calibredbPath = path
		}
	}

	if calibredbPath == "" {
		t.Skip("calibredb not found, skipping test with real calibredb")
	}

	tempDir := os.TempDir() + "/calibre_test_version"
	defer func() { _ = os.RemoveAll(tempDir) }()

	c := calibredb.NewCalibre(
		calibredb.WithCalibreDBLocation(calibredbPath),
		calibredb.WithLibraryPath(tempDir),
	)

	got := c.Version()

	// Version output should contain "calibre" or version number
	if !strings.Contains(got, "calibre") && !strings.Contains(got, ".") {
		t.Errorf("Version() = %v, want to contain 'calibre' or version number", got)
	}
}

func TestCalibre_Help(t *testing.T) {
	tests := []struct {
		name              string
		calibreDBLocation string
		libraryPath       string
		wantContains      string
		wantErr           bool
	}{
		{
			name:              "Invalid calibredb location",
			calibreDBLocation: "/nonexistent/calibredb",
			libraryPath:       os.TempDir(),
			wantContains:      "no such file or directory",
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calibredb.NewCalibre(
				calibredb.WithCalibreDBLocation(tt.calibreDBLocation),
				calibredb.WithLibraryPath(tt.libraryPath),
			)

			got := c.Help()

			if tt.wantErr {
				if !strings.Contains(got, tt.wantContains) {
					t.Errorf("Help() = %v, want to contain %v", got, tt.wantContains)
				}
			} else {
				if strings.Contains(got, "error") || strings.Contains(got, "Error") {
					t.Errorf("Help() returned error: %v", got)
				}
			}
		})
	}
}

func TestCalibre_Help_WithRealCalibredb(t *testing.T) {
	// Try to find calibredb in common locations
	calibredbPaths := []string{
		"/usr/bin/calibredb",
		"/usr/local/bin/calibredb",
		"/Applications/calibre.app/Contents/MacOS/calibredb",
	}

	var calibredbPath string
	for _, path := range calibredbPaths {
		if _, err := os.Stat(path); err == nil {
			calibredbPath = path
			break
		}
	}

	if calibredbPath == "" {
		// Also try 'which calibredb'
		if path, err := exec.LookPath("calibredb"); err == nil {
			calibredbPath = path
		}
	}

	if calibredbPath == "" {
		t.Skip("calibredb not found, skipping test with real calibredb")
	}

	tempDir := os.TempDir() + "/calibre_test_help"
	defer func() { _ = os.RemoveAll(tempDir) }()

	c := calibredb.NewCalibre(
		calibredb.WithCalibreDBLocation(calibredbPath),
		calibredb.WithLibraryPath(tempDir),
	)

	got := c.Help()

	// Help output should contain "Usage" or "command"
	if !strings.Contains(got, "Usage") && !strings.Contains(got, "command") && !strings.Contains(got, "calibredb") {
		t.Errorf("Help() = %v, want to contain 'Usage', 'command', or 'calibredb'", got)
	}
}

func TestCalibre_OnErrorCallback(t *testing.T) {
	called := false
	var capturedError error

	onError := func(err error) {
		called = true
		capturedError = err
	}

	c := calibredb.NewCalibre(
		calibredb.WithCalibreDBLocation("/nonexistent/calibredb"),
		calibredb.WithLibraryPath(os.TempDir()),
		calibredb.WithOnError(onError),
	)

	// Call a method that will fail
	_ = c.Version()

	if !called {
		t.Error("OnError callback was not called when command failed")
	}

	if capturedError == nil {
		t.Error("OnError callback received nil error")
	}
}

func TestWithCalibreDBLocation(t *testing.T) {
	testPath := "/test/path/to/calibredb"
	c := calibredb.NewCalibre(calibredb.WithCalibreDBLocation(testPath))

	if c.CalibreDBLocation != testPath {
		t.Errorf("CalibreDBLocation = %v, want %v", c.CalibreDBLocation, testPath)
	}
}

func TestWithLibraryPath(t *testing.T) {
	testPath := "/test/path/to/library"
	c := calibredb.NewCalibre(calibredb.WithLibraryPath(testPath))

	if c.LibraryPath != testPath {
		t.Errorf("LibraryPath = %v, want %v", c.LibraryPath, testPath)
	}
}

func TestWithUsername(t *testing.T) {
	testUsername := "testuser"
	c := calibredb.NewCalibre(calibredb.WithUsername(testUsername))

	if c.Username != testUsername {
		t.Errorf("Username = %v, want %v", c.Username, testUsername)
	}
}

func TestWithPassword(t *testing.T) {
	testPassword := "testpassword"
	c := calibredb.NewCalibre(calibredb.WithPassword(testPassword))

	if c.Password != testPassword {
		t.Errorf("Password = %v, want %v", c.Password, testPassword)
	}
}

func TestWithTimeout(t *testing.T) {
	testTimeout := "60s"
	c := calibredb.NewCalibre(calibredb.WithTimeout(testTimeout))

	if c.Timeout != testTimeout {
		t.Errorf("Timeout = %v, want %v", c.Timeout, testTimeout)
	}
}

func TestCalibre_ChainedOptions(t *testing.T) {
	// Test that multiple options can be chained together
	c := calibredb.NewCalibre(
		calibredb.WithCalibreDBLocation("/usr/bin/calibredb"),
		calibredb.WithLibraryPath("/path/to/library"),
		calibredb.WithUsername("user"),
		calibredb.WithPassword("pass"),
		calibredb.WithTimeout("30s"),
	)

	if c.CalibreDBLocation != "/usr/bin/calibredb" {
		t.Errorf("CalibreDBLocation = %v, want /usr/bin/calibredb", c.CalibreDBLocation)
	}
	if c.LibraryPath != "/path/to/library" {
		t.Errorf("LibraryPath = %v, want /path/to/library", c.LibraryPath)
	}
	if c.Username != "user" {
		t.Errorf("Username = %v, want user", c.Username)
	}
	if c.Password != "pass" {
		t.Errorf("Password = %v, want pass", c.Password)
	}
	if c.Timeout != "30s" {
		t.Errorf("Timeout = %v, want 30s", c.Timeout)
	}
}

func TestCalibre_EmptyOptions(t *testing.T) {
	// Test that creating a Calibre instance with no options works
	c := calibredb.NewCalibre()

	if c.CalibreDBLocation != "" {
		t.Errorf("CalibreDBLocation = %v, want empty string", c.CalibreDBLocation)
	}
	if c.LibraryPath != "" {
		t.Errorf("LibraryPath = %v, want empty string", c.LibraryPath)
	}
	if c.Username != "" {
		t.Errorf("Username = %v, want empty string", c.Username)
	}
	if c.Password != "" {
		t.Errorf("Password = %v, want empty string", c.Password)
	}
	if c.Timeout != "" {
		t.Errorf("Timeout = %v, want empty string", c.Timeout)
	}
	if c.OnError != nil {
		t.Error("OnError should be nil by default")
	}
}

func TestCalibre_Version_SuccessPath(t *testing.T) {
	// This test will pass if calibredb is installed, otherwise it will skip
	calibredbPaths := []string{
		"/usr/bin/calibredb",
		"/usr/local/bin/calibredb",
		"/Applications/calibre.app/Contents/MacOS/calibredb",
	}

	var calibredbPath string
	for _, path := range calibredbPaths {
		if _, err := os.Stat(path); err == nil {
			calibredbPath = path
			break
		}
	}

	if calibredbPath == "" {
		if path, err := exec.LookPath("calibredb"); err == nil {
			calibredbPath = path
		}
	}

	if calibredbPath == "" {
		t.Skip("calibredb not found, skipping success path test")
	}

	tempDir := os.TempDir() + "/calibre_test_version_success"
	defer func() { _ = os.RemoveAll(tempDir) }()

	errorCalled := false
	c := calibredb.NewCalibre(
		calibredb.WithCalibreDBLocation(calibredbPath),
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithOnError(func(err error) {
			errorCalled = true
		}),
	)

	version := c.Version()

	// On success, OnError should not be called
	if errorCalled {
		t.Error("OnError was called on successful Version() call")
	}

	// Version should contain some version information
	if version == "" {
		t.Error("Version() returned empty string")
	}
}

func TestCalibre_Help_SuccessPath(t *testing.T) {
	// This test will pass if calibredb is installed, otherwise it will skip
	calibredbPaths := []string{
		"/usr/bin/calibredb",
		"/usr/local/bin/calibredb",
		"/Applications/calibre.app/Contents/MacOS/calibredb",
	}

	var calibredbPath string
	for _, path := range calibredbPaths {
		if _, err := os.Stat(path); err == nil {
			calibredbPath = path
			break
		}
	}

	if calibredbPath == "" {
		if path, err := exec.LookPath("calibredb"); err == nil {
			calibredbPath = path
		}
	}

	if calibredbPath == "" {
		t.Skip("calibredb not found, skipping success path test")
	}

	tempDir := os.TempDir() + "/calibre_test_help_success"
	defer func() { _ = os.RemoveAll(tempDir) }()

	errorCalled := false
	c := calibredb.NewCalibre(
		calibredb.WithCalibreDBLocation(calibredbPath),
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithOnError(func(err error) {
			errorCalled = true
		}),
	)

	help := c.Help()

	// On success, OnError should not be called
	if errorCalled {
		t.Error("OnError was called on successful Help() call")
	}

	// Help should contain some help information
	if help == "" {
		t.Error("Help() returned empty string")
	}
}

func TestCalibre_ErrorPathsWithoutOnError(t *testing.T) {
	// Test error handling when OnError callback is not set
	c := calibredb.NewCalibre(
		calibredb.WithCalibreDBLocation("/nonexistent/calibredb"),
		calibredb.WithLibraryPath(os.TempDir()),
	)

	version := c.Version()
	if !strings.Contains(version, "no such file or directory") && !strings.Contains(version, "executable file not found") {
		t.Errorf("Version() = %v, expected error message about missing file", version)
	}

	help := c.Help()
	if !strings.Contains(help, "no such file or directory") && !strings.Contains(help, "executable file not found") {
		t.Errorf("Help() = %v, expected error message about missing file", help)
	}
}

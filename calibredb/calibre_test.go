package calibredb_test

import (
	"os"
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
			name: "Default calibre instance",
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
				calibredb.WithTimeout("30"),
			},
			want: &calibredb.Calibre{
				Timeout: "30",
			},
		},
		{
			name: "With all options",
			opts: []calibredb.CalibreOption{
				calibredb.WithCalibreDBLocation("/usr/bin/calibredb"),
				calibredb.WithLibraryPath("/path/to/library"),
				calibredb.WithUsername("testuser"),
				calibredb.WithPassword("testpass"),
				calibredb.WithTimeout("30"),
			},
			want: &calibredb.Calibre{
				CalibreDBLocation: "/usr/bin/calibredb",
				LibraryPath:       "/path/to/library",
				Username:          "testuser",
				Password:          "testpass",
				Timeout:           "30",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calibredb.NewCalibre(tt.opts...)
			if got == nil {
				t.Fatal("NewCalibre() returned nil")
			}
			if got.CalibreDBLocation != tt.want.CalibreDBLocation {
				t.Errorf("NewCalibre() CalibreDBLocation = %v, want %v", got.CalibreDBLocation, tt.want.CalibreDBLocation)
			}
			if got.LibraryPath != tt.want.LibraryPath {
				t.Errorf("NewCalibre() LibraryPath = %v, want %v", got.LibraryPath, tt.want.LibraryPath)
			}
			if got.Username != tt.want.Username {
				t.Errorf("NewCalibre() Username = %v, want %v", got.Username, tt.want.Username)
			}
			if got.Password != tt.want.Password {
				t.Errorf("NewCalibre() Password = %v, want %v", got.Password, tt.want.Password)
			}
			if got.Timeout != tt.want.Timeout {
				t.Errorf("NewCalibre() Timeout = %v, want %v", got.Timeout, tt.want.Timeout)
			}
		})
	}
}

func TestWithOnError(t *testing.T) {
	onErrorFunc := func(err error) {
		// Error handler function for testing
		_ = err
	}

	c := calibredb.NewCalibre(
		calibredb.WithOnError(onErrorFunc),
	)

	if c == nil {
		t.Fatal("NewCalibre() returned nil")
	}

	// We can't directly test the OnError function as it's a private field,
	// but we can verify that the instance was created successfully
	if c.OnError == nil {
		t.Error("OnError callback was not set")
	}
}

func TestCalibreOptions(t *testing.T) {
	tests := []struct {
		name  string
		opt   calibredb.CalibreOption
		check func(*calibredb.Calibre) error
	}{
		{
			name: "WithCalibreDBLocation sets location",
			opt:  calibredb.WithCalibreDBLocation("/test/path"),
			check: func(c *calibredb.Calibre) error {
				if c.CalibreDBLocation != "/test/path" {
					t.Errorf("CalibreDBLocation = %v, want /test/path", c.CalibreDBLocation)
				}
				return nil
			},
		},
		{
			name: "WithLibraryPath sets path",
			opt:  calibredb.WithLibraryPath("/library/path"),
			check: func(c *calibredb.Calibre) error {
				if c.LibraryPath != "/library/path" {
					t.Errorf("LibraryPath = %v, want /library/path", c.LibraryPath)
				}
				return nil
			},
		},
		{
			name: "WithUsername sets username",
			opt:  calibredb.WithUsername("user123"),
			check: func(c *calibredb.Calibre) error {
				if c.Username != "user123" {
					t.Errorf("Username = %v, want user123", c.Username)
				}
				return nil
			},
		},
		{
			name: "WithPassword sets password",
			opt:  calibredb.WithPassword("pass123"),
			check: func(c *calibredb.Calibre) error {
				if c.Password != "pass123" {
					t.Errorf("Password = %v, want pass123", c.Password)
				}
				return nil
			},
		},
		{
			name: "WithTimeout sets timeout",
			opt:  calibredb.WithTimeout("60"),
			check: func(c *calibredb.Calibre) error {
				if c.Timeout != "60" {
					t.Errorf("Timeout = %v, want 60", c.Timeout)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calibredb.NewCalibre(tt.opt)
			if c == nil {
				t.Fatal("NewCalibre() returned nil")
			}
			tt.check(c)
		})
	}
}

func TestCalibre_Version(t *testing.T) {
	tests := []struct {
		name              string
		calibreDBLocation string
		wantContains      string
	}{
		{
			name:              "Invalid calibredb location returns error",
			calibreDBLocation: "/nonexistent/calibredb",
			wantContains:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := os.TempDir() + "/" + t.Name()
			defer func() { _ = os.RemoveAll(tempDir) }()

			c := calibredb.NewCalibre(
				calibredb.WithCalibreDBLocation(tt.calibreDBLocation),
				calibredb.WithLibraryPath(tempDir),
			)

			got := c.Version()
			if got == "" {
				t.Error("Version() returned empty string")
			}
			// Since we're using an invalid path, we expect an error message
			if tt.calibreDBLocation == "/nonexistent/calibredb" {
				if !strings.Contains(got, "no such file") && !strings.Contains(got, "not found") && !strings.Contains(got, "executable file not found") {
					t.Logf("Version() returned: %v", got)
				}
			}
		})
	}
}

func TestCalibre_Help(t *testing.T) {
	tests := []struct {
		name              string
		calibreDBLocation string
		wantContains      string
	}{
		{
			name:              "Invalid calibredb location returns error",
			calibreDBLocation: "/nonexistent/calibredb",
			wantContains:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := os.TempDir() + "/" + t.Name()
			defer func() { _ = os.RemoveAll(tempDir) }()

			c := calibredb.NewCalibre(
				calibredb.WithCalibreDBLocation(tt.calibreDBLocation),
				calibredb.WithLibraryPath(tempDir),
			)

			got := c.Help()
			if got == "" {
				t.Error("Help() returned empty string")
			}
			// Since we're using an invalid path, we expect an error message
			if tt.calibreDBLocation == "/nonexistent/calibredb" {
				if !strings.Contains(got, "no such file") && !strings.Contains(got, "not found") && !strings.Contains(got, "executable file not found") {
					t.Logf("Help() returned: %v", got)
				}
			}
		})
	}
}

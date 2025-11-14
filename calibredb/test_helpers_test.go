package calibredb_test

import "os"

// getCalibreDBPath returns the path to calibredb binary.
// It first checks the CALIBREDB_PATH environment variable,
// then falls back to the macOS default path.
func getCalibreDBPath() string {
	if path := os.Getenv("CALIBREDB_PATH"); path != "" {
		return path
	}
	// Default to macOS path for local development
	return "/Applications/calibre.app/Contents/MacOS/calibredb"
}

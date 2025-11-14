package calibredb_test

import (
	"os"

	"github.com/veverkap/calibre-rest/calibredb"
)

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

func getTestCalibre(name string) (*calibredb.Calibre, func()) {
	tempDir := os.TempDir() + "/" + name

	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation(getCalibreDBPath()),
	)
	return c, func() {
		_ = os.RemoveAll(tempDir)
	}
}

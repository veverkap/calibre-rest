package calibredb_test

import (
	"os"

	"github.com/veverkap/calibre-rest/calibredb"
)

func getTestCalibre(name string) (*calibredb.Calibre, func()) {
	path := os.Getenv("CALIBREDB_PATH")
	if path == "" {
		path = "/Applications/calibre.app/Contents/MacOS/calibredb"
	}
	tempDir := os.TempDir() + "/" + name

	c := calibredb.NewCalibre(
		calibredb.WithLibraryPath(tempDir),
		calibredb.WithCalibreDBLocation(path),
	)
	return c, func() {
		_ = os.RemoveAll(tempDir)
	}
}

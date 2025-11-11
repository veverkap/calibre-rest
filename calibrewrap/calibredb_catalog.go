package calibrewrap
	
// CalibredbCATALOG executes `calibredb catalog`.
//
//
// Flags:
//   --ids: Comma-separated list of database IDs to catalog.
//   --search: ignored.
//   --verbose: Show detailed output information. Useful for debugging
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre
//   --generate-authors: Include 'Authors' section in catalog.
//   --generate-genres: Include 'Genres' section in catalog.
//   --generate-titles: Include 'Titles' section in catalog.
//   --generate-series: Include 'Series' section in catalog.
//   --preset: Use a named preset created with the GUI catalog

type CalibredbCATALOGFlags struct {
	Ids string // Comma-separated list of database IDs to catalog.
	Search string // ignored.
	Verbose string // Show detailed output information. Useful for debugging
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
	GenerateAuthors string // Include 'Authors' section in catalog.
	GenerateGenres string // Include 'Genres' section in catalog.
	GenerateTitles string // Include 'Titles' section in catalog.
	GenerateSeries string // Include 'Series' section in catalog.
	Preset string // Use a named preset created with the GUI catalog
}

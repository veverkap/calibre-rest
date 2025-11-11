package calibrewrap
	
// CalibredbFTS_SEARCH executes `calibredb fts_search`.
//
//
// Flags:
//   --include-snippets: Include snippets of the text surrounding each match.
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre

type CalibredbFTS_SEARCHFlags struct {
	IncludeSnippets string // Include snippets of the text surrounding each match.
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
}

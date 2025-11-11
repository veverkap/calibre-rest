package calibrewrap
	
// CalibredbREMOVE executes `calibredb remove`.
//
// Remove a book from the calibre database.
//
// Flags:
//   --permanent: Do not use the Trash
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre

type CalibredbREMOVEFlags struct {
	Permanent string // Do not use the Trash
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
}

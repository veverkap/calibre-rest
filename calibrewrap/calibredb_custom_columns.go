package calibrewrap
	
// CalibredbCUSTOM_COLUMNS executes `calibredb custom_columns`.
//
//
// Flags:
//   --details: Show details for each column.
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre

type CalibredbCUSTOM_COLUMNSFlags struct {
	Details string // Show details for each column.
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
}

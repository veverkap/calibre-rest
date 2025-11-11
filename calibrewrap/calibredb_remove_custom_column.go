package calibrewrap
	
// CalibredbREMOVE_CUSTOM_COLUMN executes `calibredb remove_custom_column`.
//
//
// Flags:
//   --force: Do not ask for confirmation
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre

type CalibredbREMOVE_CUSTOM_COLUMNFlags struct {
	Force string // Do not ask for confirmation
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
}

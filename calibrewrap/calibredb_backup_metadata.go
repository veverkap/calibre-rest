package calibrewrap
	
// CalibredbBACKUP_METADATA executes `calibredb backup_metadata`.
//
//
// Flags:
//   --all: option.
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre

type CalibredbBACKUP_METADATAFlags struct {
	All string // option.
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
}

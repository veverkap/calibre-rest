package calibrewrap
	
// CalibredbSET_METADATA executes `calibredb set_metadata`.
//
//
// Flags:
//   --field: If you use the --field option, there
//   --list-fields: List the metadata field names that can be used with
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre

type CalibredbSET_METADATAFlags struct {
	Field string // If you use the --field option, there
	ListFields string // List the metadata field names that can be used with
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
}

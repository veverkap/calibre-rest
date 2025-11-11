package calibrewrap
	
// CalibredbSHOW_METADATA executes `calibredb show_metadata`.
//
//
// Flags:
//   --as-opf: Print metadata in OPF form (XML)
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre

type CalibredbSHOW_METADATAFlags struct {
	AsOpf string // Print metadata in OPF form (XML)
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
}

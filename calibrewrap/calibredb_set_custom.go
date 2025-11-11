package calibrewrap
	
// CalibredbSET_CUSTOM executes `calibredb set_custom`.
//
//
// Flags:
//   --append: If the column stores multiple values, append the
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre

type CalibredbSET_CUSTOMFlags struct {
	Append string // If the column stores multiple values, append the
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
}

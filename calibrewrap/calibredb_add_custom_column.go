package calibrewrap
	
// CalibredbADD_CUSTOM_COLUMN executes `calibredb add_custom_column`.
//
//
// Flags:
//   --is-multiple: This column stores tag like data (i.e. multiple comma
//   --display: A dictionary of options to customize how the data in
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre

type CalibredbADD_CUSTOM_COLUMNFlags struct {
	IsMultiple string // This column stores tag like data (i.e. multiple comma
	Display string // A dictionary of options to customize how the data in
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
}

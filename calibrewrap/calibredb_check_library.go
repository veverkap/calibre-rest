package calibrewrap
	
// CalibredbCHECK_LIBRARY executes `calibredb check_library`.
//
//
// Flags:
//   --csv: Output in CSV
//   --vacuum-fts-db: Vacuum the full text search database. This can be very
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre

type CalibredbCHECK_LIBRARYFlags struct {
	Csv string // Output in CSV
	VacuumFtsDb string // Vacuum the full text search database. This can be very
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
}

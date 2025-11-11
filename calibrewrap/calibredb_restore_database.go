package calibrewrap
	
// CalibredbRESTORE_DATABASE executes `calibredb restore_database`.
//
//
// Flags:
//   --really-do-it: Really do the recovery. The command will not run
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre

type CalibredbRESTORE_DATABASEFlags struct {
	ReallyDoIt string // Really do the recovery. The command will not run
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
}

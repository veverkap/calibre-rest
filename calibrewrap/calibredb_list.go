package calibrewrap
	
// CalibredbLIST executes `calibredb list`.
//
// List all books in the calibre database.
//
// Flags:
//   --sort-by: The field by which to sort the results. You can
//   --ascending: Sort results in ascending order
//   --prefix: The prefix for all file paths. Default is the absolute
//   --limit: The maximum number of results to display. Default: all
//   --for-machine: Generate output in JSON format, which is more suitable
//   --template: The template to run if "template" is in the field
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre

type CalibredbLISTFlags struct {
	SortBy string // The field by which to sort the results. You can
	Ascending string // Sort results in ascending order
	Prefix string // The prefix for all file paths. Default is the absolute
	Limit string // The maximum number of results to display. Default: all
	ForMachine string // Generate output in JSON format, which is more suitable
	Template string // The template to run if "template" is in the field
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
}

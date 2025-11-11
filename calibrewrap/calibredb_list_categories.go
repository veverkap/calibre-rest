package calibrewrap
	
// CalibredbLIST_CATEGORIES executes `calibredb list_categories`.
//
//
// Flags:
//   --item_count: Output only the number of items in a category instead
//   --csv: Output in CSV
//   --dialect: The type of CSV file to produce. Choices: excel,
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre

type CalibredbLIST_CATEGORIESFlags struct {
	Item_count string // Output only the number of items in a category instead
	Csv string // Output in CSV
	Dialect string // The type of CSV file to produce. Choices: excel,
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
}

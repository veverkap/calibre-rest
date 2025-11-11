package calibrewrap
	
// CalibredbADD executes `calibredb add`.
//
// Add a book to the calibre database.
//
// Flags:
//   --duplicates: Add books to database even if they already exist.
//   --automerge: takes precedence.
//   --empty: Add an empty book (a book with no formats)
//   --isbn: Set the ISBN of the added book(s)
//   --tags: Set the tags of the added book(s)
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre
//   --recurse: Process folders recursively
//   --ignore: PATTERN
//   --add: PATTERN  A filename (glob) pattern, files matching this pattern

type CalibredbADDFlags struct {
	Duplicates string // Add books to database even if they already exist.
	Automerge string // takes precedence.
	Empty string // Add an empty book (a book with no formats)
	Isbn string // Set the ISBN of the added book(s)
	Tags string // Set the tags of the added book(s)
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
	Recurse string // Process folders recursively
	Ignore string // PATTERN
	Add string // PATTERN  A filename (glob) pattern, files matching this pattern
}

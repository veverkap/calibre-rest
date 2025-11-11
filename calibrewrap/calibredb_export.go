package calibrewrap
	
// CalibredbEXPORT executes `calibredb export`.
//
//
// Flags:
//   --all: Export all books in database, ignoring the list of
//   --to-dir: Export books to the specified folder. Default is .
//   --single-dir: Export all books into a single folder
//   --progress: Report progress
//   --dont-asciiize: Have calibre convert all non English characters into
//   --dont-write-opf: Normally, calibre will write the metadata into a
//   --dont-save-cover: Normally, calibre will save the cover in a separate
//   --timefmt: The format in which to display dates. %d - day, %b -
//   --template: The template to control the filename and folder
//   --formats: Comma separated list of formats to save for each book.
//   --replace-whitespace: Replace whitespace with underscores.
//   --to-lowercase: Convert paths to lowercase.
//   --library-path: --with-library=LIBRARY_PATH
//   --help: show this help message and exit
//   --version: show program's version number and exit
//   --timeout: The timeout, in seconds, when connecting to a calibre

type CalibredbEXPORTFlags struct {
	All string // Export all books in database, ignoring the list of
	ToDir string // Export books to the specified folder. Default is .
	SingleDir string // Export all books into a single folder
	Progress string // Report progress
	DontAsciiize string // Have calibre convert all non English characters into
	DontWriteOpf string // Normally, calibre will write the metadata into a
	DontSaveCover string // Normally, calibre will save the cover in a separate
	Timefmt string // The format in which to display dates. %d - day, %b -
	Template string // The template to control the filename and folder
	Formats string // Comma separated list of formats to save for each book.
	ReplaceWhitespace string // Replace whitespace with underscores.
	ToLowercase string // Convert paths to lowercase.
	LibraryPath string // --with-library=LIBRARY_PATH
	Help string // show this help message and exit
	Version string // show program's version number and exit
	Timeout string // The timeout, in seconds, when connecting to a calibre
}

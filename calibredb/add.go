package calibredb
	// Usage: calibredb add [options] file1 file2 file3 ...
// 
// Add the specified files as books to the database. You can also specify folders, see
// the folder related options below.
// 
// 
// Whenever you pass arguments to calibredb that have spaces in them, enclose the arguments in quotation marks. For example: "/some path/with spaces"
// 
// Options:
//   -d, --duplicates      Add books to database even if they already exist.
//                         Comparison is done based on book titles and authors.
//                         Note that the --automerge option takes precedence.
// 
//   -m AUTOMERGE, --automerge=AUTOMERGE
//                         If books with similar titles and authors are found,
//                         merge the incoming formats (files) automatically into
//                         existing book records. A value of "ignore" means
//                         duplicate formats are discarded. A value of
//                         "overwrite" means duplicate formats in the library are
//                         overwritten with the newly added files. A value of
//                         "new_record" means duplicate formats are placed into a
//                         new book record.
// 
//   -e, --empty           Add an empty book (a book with no formats)
// 
//   -t TITLE, --title=TITLE
//                         Set the title of the added book(s)
// 
//   -a AUTHORS, --authors=AUTHORS
//                         Set the authors of the added book(s)
// 
//   -i ISBN, --isbn=ISBN  Set the ISBN of the added book(s)
// 
//   -I IDENTIFIER, --identifier=IDENTIFIER
//                         Set the identifiers for this book, e.g. -I asin:XXX -I
//                         isbn:YYY
// 
//   -T TAGS, --tags=TAGS  Set the tags of the added book(s)
// 
//   -s SERIES, --series=SERIES
//                         Set the series of the added book(s)
// 
//   -S SERIES_INDEX, --series-index=SERIES_INDEX
//                         Set the series number of the added book(s)
// 
//   -c COVER, --cover=COVER
//                         Path to the cover to use for the added book
// 
//   -l LANGUAGES, --languages=LANGUAGES
//                         A comma separated list of languages (best to use
//                         ISO639 language codes, though some language names may
//                         also be recognized)
// 
// 
//   GLOBAL OPTIONS:
//     --library-path=LIBRARY_PATH, --with-library=LIBRARY_PATH
//                         Path to the calibre library. Default is to use the
//                         path stored in the settings. You can also connect to a
//                         calibre Content server to perform actions on remote
//                         libraries. To do so use a URL of the form:
//                         http://hostname:port/#library_id for example,
//                         http://localhost:8080/#mylibrary. library_id is the
//                         library id of the library you want to connect to on
//                         the Content server. You can use the special library_id
//                         value of - to get a list of library ids available on
//                         the server. For details on how to setup access via a
//                         Content server, see https://manual.calibre-
//                         ebook.com/generated/en/calibredb.html.
// 
//     -h, --help          show this help message and exit
// 
//     --version           show program's version number and exit
// 
//     --username=USERNAME
//                         Username for connecting to a calibre Content server
// 
//     --password=PASSWORD
//                         Password for connecting to a calibre Content server.
//                         To read the password from standard input, use the
//                         special value: <stdin>. To read the password from a
//                         file, use: <f:/path/to/file> (i.e. <f: followed by the
//                         full path to the file and a trailing >). The angle
//                         brackets in the above are required, remember to escape
//                         them or use quotes for your shell.
// 
//     --timeout=TIMEOUT   The timeout, in seconds, when connecting to a calibre
//                         library over the network. The default is two minutes.
// 
// 
//   ADDING FROM FOLDERS:
//     Options to control the adding of books from folders. By default only
//     files that have extensions of known e-book file types are added.
// 
//     -1, --one-book-per-directory
//                         Assume that each folder has only a single logical book
//                         and that all files in it are different e-book formats
//                         of that book
// 
//     -r, --recurse       Process folders recursively
// 
//     --ignore=GLOB PATTERN
//                         A filename (glob) pattern, files matching this pattern
//                         will be ignored when scanning folders for files. Can
//                         be specified multiple times for multiple patterns. For
//                         example: *.pdf will ignore all PDF files
// 
//     --add=GLOB PATTERN  A filename (glob) pattern, files matching this pattern
//                         will be added when scanning folders for files, even if
//                         they are not of a known e-book file type. Can be
//                         specified multiple times for multiple patterns.
// 
// 
// Created by Kovid Goyal <kovid@kovidgoyal.net>
// 
type AddOptions struct {
	Empty string `json:"empty,omitempty"`
	Isbn string `json:"isbn,omitempty"`
	Tags string `json:"tags,omitempty"`
	Recurse string `json:"recurse,omitempty"`
	Ignore string `json:"ignore,omitempty"`
	Add string `json:"add,omitempty"`
	Duplicates string `json:"duplicates,omitempty"`
	Automerge string `json:"automerge,omitempty"`
}

func (c *Calibre) Add(opts AddOptions, args ...string) string {
	return "add"
}

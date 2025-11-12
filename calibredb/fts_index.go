package calibredb
	// Usage: calibredb fts_index [options] enable/disable/status/reindex
// 
// Control the Full text search indexing process.
// 
// enable
//     Turns on FTS indexing for this library
// disable
//     Turns off FTS indexing for this library
// status
//     Shows the current indexing status
// reindex
//     Can be used to re-index either particular books or
//     the entire library. To re-index particular books
//     specify the book ids as additional arguments after the
//     reindex command. If no book ids are specified the
//     entire library is re-indexed.
// 
// 
// Whenever you pass arguments to calibredb that have spaces in them, enclose the arguments in quotation marks. For example: "/some path/with spaces"
// 
// Options:
//   --wait-for-completion
//                         Wait till all books are indexed, showing indexing
//                         progress periodically
// 
//   --indexing-speed=INDEXING_SPEED
//                         The speed of indexing. Use fast for fast indexing
//                         using all your computers resources and slow for less
//                         resource intensive indexing. Note that the speed is
//                         reset to slow after every invocation.
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
// Created by Kovid Goyal <kovid@kovidgoyal.net>
// 
type FtsIndexOptions struct {
}

func (c *Calibre) FtsIndex(opts FtsIndexOptions, args ...string) string {
	return "fts_index"
}

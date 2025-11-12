package calibredb
		
	// Usage: calibredb add_format [options] id ebook_file
// 
// Add the e-book in ebook_file to the available formats for the logical book identified by id. You can get id by using the search command. If the format already exists, it is replaced, unless the do not replace option is specified.
// 
// Whenever you pass arguments to calibredb that have spaces in them, enclose the arguments in quotation marks. For example: "/some path/with spaces"
// 
// Options:
//   --dont-replace        Do not replace the format if it already exists
// 
//   --as-extra-data-file  Add the file as an extra data file to the book, not an
//                         ebook format
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
type AddFormatOptions struct {
	AsExtraDataFile string `json:"as-extra-data-file,omitempty"`
	DontReplace string `json:"dont-replace,omitempty"`
}

func (c *Calibre) AddFormatHelp() string {
	return c.run("add_format", "-h")
}

func (c *Calibre) AddFormat(opts AddFormatOptions, args ...string) string {
	return "add_format"
}

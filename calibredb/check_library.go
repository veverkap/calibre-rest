package calibredb
		
	// Usage: calibredb check_library [options]
// 
// Perform some checks on the filesystem representing a library. Reports are invalid_titles, extra_titles, invalid_authors, extra_authors, missing_formats, extra_formats, extra_files, missing_covers, extra_covers, malformed_formats, malformed_paths, failed_folders
// 
// 
// Whenever you pass arguments to calibredb that have spaces in them, enclose the arguments in quotation marks. For example: "/some path/with spaces"
// 
// Options:
//   -c, --csv             Output in CSV
// 
//   -r REPORT, --report=REPORT
//                         Comma-separated list of reports.
//                         Default: all
// 
//   -e EXTS, --ignore_extensions=EXTS
//                         Comma-separated list of extensions to ignore.
//                         Default: all
// 
//   -n NAMES, --ignore_names=NAMES
//                         Comma-separated list of names to ignore.
//                         Default: all
// 
//   --vacuum-fts-db       Vacuum the full text search database. This can be very
//                         slow and memory intensive, depending on the size of
//                         the database.
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
type CheckLibraryOptions struct {
	Csv string `json:"csv,omitempty"`
	VacuumFtsDb string `json:"vacuum-fts-db,omitempty"`
}

func (c *Calibre) CheckLibraryHelp() string {
	return c.run("check_library", "-h")
}

func (c *Calibre) CheckLibrary(opts CheckLibraryOptions, args ...string) string {
	return "check_library"
}

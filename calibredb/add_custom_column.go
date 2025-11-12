package calibredb
	// Usage: calibredb add_custom_column [options] label name datatype
// 
// Create a custom column. label is the machine friendly name of the column. Should
// not contain spaces or colons. name is the human friendly name of the column.
// datatype is one of: bool, comments, composite, datetime, enumeration, float, int, rating, series, text
// 
// 
// Whenever you pass arguments to calibredb that have spaces in them, enclose the arguments in quotation marks. For example: "/some path/with spaces"
// 
// Options:
//   --is-multiple         This column stores tag like data (i.e. multiple comma
//                         separated values). Only applies if datatype is text.
// 
//   --display=DISPLAY     A dictionary of options to customize how the data in
//                         this column will be interpreted. This is a JSON
//                         string. For enumeration columns, use
//                         --display="{\"enum_values\":[\"val1\", \"val2\"]}"
//                         There are many options that can go into the display
//                         variable.The options by column type are:
//                         composite: composite_template, composite_sort,
//                         make_category,contains_html, use_decorations
//                         datetime: date_format
//                         enumeration: enum_values, enum_colors, use_decorations
//                         int, float: number_format
//                         text: is_names, use_decorations
//                         The best way to find legal combinations is to create a
//                         custom column of the appropriate type in the GUI then
//                         look at the backup OPF for a book (ensure that a new
//                         OPF has been created since the column was added). You
//                         will see the JSON for the "display" for the new column
//                         in the OPF.
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
type AddCustomColumnOptions struct {
	IsMultiple string `json:"is-multiple,omitempty"`
	Display string `json:"display,omitempty"`
}

func (c *Calibre) AddCustomColumn(opts AddCustomColumnOptions, args ...string) string {
	return "add_custom_column"
}

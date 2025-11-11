Integration status: False
Usage: calibredb list [options]

List the books available in the calibre database.


Whenever you pass arguments to calibredb that have spaces in them, enclose the arguments in quotation marks. For example: "/some path/with spaces"

Options:
  -f FIELDS, --fields=FIELDS
                        The fields to display when listing books in the
                        database. Should be a comma separated list of fields.
                        Available fields: author_sort, authors, comments,
                        cover, formats, identifiers, isbn, languages,
                        last_modified, pubdate, publisher, rating, series,
                        series_index, size, tags, template, timestamp, title,
                        uuid
                        Default: title,authors. The special field "all" can be
                        used to select all fields. In addition to the builtin
                        fields above, custom fields are also available as
                        *field_name, for example, for a custom field #rating,
                        use the name: *rating

  --sort-by=SORT_BY     The field by which to sort the results. You can
                        specify multiple fields by separating them with
                        commas.
                        Available fields: author_sort, authors, comments,
                        cover, formats, identifiers, isbn, languages,
                        last_modified, pubdate, publisher, rating, series,
                        series_index, size, tags, template, timestamp, title,
                        uuid
                        Default: id. In addition to the builtin fields above,
                        custom fields are also available as *field_name, for
                        example, for a custom field #rating, use the name:
                        *rating

  --ascending           Sort results in ascending order

  -s SEARCH, --search=SEARCH
                        Filter the results by the search query. For the format
                        of the search query, please see the search related
                        documentation in the User Manual. Default is to do no
                        filtering.

  -w LINE_WIDTH, --line-width=LINE_WIDTH
                        The maximum width of a single line in the output.
                        Defaults to detecting screen size.

  --separator=SEPARATOR
                        The string used to separate fields. Default is a
                        space.

  --prefix=PREFIX       The prefix for all file paths. Default is the absolute
                        path to the library folder.

  --limit=LIMIT         The maximum number of results to display. Default: all

  --for-machine         Generate output in JSON format, which is more suitable
                        for machine parsing. Causes the line width and
                        separator options to be ignored.

  --template=TEMPLATE   The template to run if "template" is in the field
                        list. Note that templates are ignored while connecting
                        to a calibre server. Default: None

  -t TEMPLATE_FILE, --template_file=TEMPLATE_FILE
                        Path to a file containing the template to run if
                        "template" is in the field list. Default: None

  --template_heading=TEMPLATE_HEADING
                        Heading for the template column. Default: template.
                        This option is ignored if the option --for-machine is
                        set


  GLOBAL OPTIONS:
    --library-path=LIBRARY_PATH, --with-library=LIBRARY_PATH
                        Path to the calibre library. Default is to use the
                        path stored in the settings. You can also connect to a
                        calibre Content server to perform actions on remote
                        libraries. To do so use a URL of the form:
                        http://hostname:port/#library_id for example,
                        http://localhost:8080/#mylibrary. library_id is the
                        library id of the library you want to connect to on
                        the Content server. You can use the special library_id
                        value of - to get a list of library ids available on
                        the server. For details on how to setup access via a
                        Content server, see https://manual.calibre-
                        ebook.com/generated/en/calibredb.html.

    -h, --help          show this help message and exit

    --version           show program's version number and exit

    --username=USERNAME
                        Username for connecting to a calibre Content server

    --password=PASSWORD
                        Password for connecting to a calibre Content server.
                        To read the password from standard input, use the
                        special value: <stdin>. To read the password from a
                        file, use: <f:/path/to/file> (i.e. <f: followed by the
                        full path to the file and a trailing >). The angle
                        brackets in the above are required, remember to escape
                        them or use quotes for your shell.

    --timeout=TIMEOUT   The timeout, in seconds, when connecting to a calibre
                        library over the network. The default is two minutes.


Created by Kovid Goyal <kovid@kovidgoyal.net>

Integration status: False
Usage: 
calibredb set_metadata [options] book_id [/path/to/metadata.opf]

Set the metadata stored in the calibre database for the book identified by
book_id from the OPF file metadata.opf. book_id is a book id number from the
search command. You can get a quick feel for the OPF format by using the
--as-opf switch to the show_metadata command. You can also set the metadata of
individual fields with the --field option. If you use the --field option, there
is no need to specify an OPF file.


Whenever you pass arguments to calibredb that have spaces in them, enclose the arguments in quotation marks. For example: "/some path/with spaces"

Options:
  -f FIELD, --field=FIELD
                        The field to set. Format is field_name:value, for
                        example: --field tags:tag1,tag2. Use --list-fields to
                        get a list of all field names. You can specify this
                        option multiple times to set multiple fields. Note:
                        For languages you must use the ISO639 language codes
                        (e.g. en for English, fr for French and so on). For
                        identifiers, the syntax is --field
                        identifiers:isbn:XXXX,doi:YYYYY. For boolean (yes/no)
                        fields use true and false or yes and no.

  -l, --list-fields     List the metadata field names that can be used with
                        the --field option


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

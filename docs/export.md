Integration status: False
Usage: calibredb export [options] ids

Export the books specified by ids (a comma separated list) to the filesystem.
The export operation saves all formats of the book, its cover and metadata (in
an OPF file). Any extra data files associated with the book are also saved.
You can get id numbers from the search command.


Whenever you pass arguments to calibredb that have spaces in them, enclose the arguments in quotation marks. For example: "/some path/with spaces"

Options:
  --all                 Export all books in database, ignoring the list of
                        ids.

  --to-dir=TO_DIR       Export books to the specified folder. Default is .

  --single-dir          Export all books into a single folder

  --progress            Report progress

  --dont-asciiize       Have calibre convert all non English characters into
                        English equivalents for the file names. This is useful
                        if saving to a legacy filesystem without full support
                        for Unicode filenames. Specifying this switch will
                        turn this behavior off.

  --dont-update-metadata
                        Normally, calibre will update the metadata in the
                        saved files from what is in the calibre library. Makes
                        saving to disk slower. Specifying this switch will
                        turn this behavior off.

  --dont-write-opf      Normally, calibre will write the metadata into a
                        separate OPF file along with the actual e-book files.
                        Specifying this switch will turn this behavior off.

  --dont-save-cover     Normally, calibre will save the cover in a separate
                        file along with the actual e-book files. Specifying
                        this switch will turn this behavior off.

  --dont-save-extra-files
                        Save any data files associated with the book when
                        saving the book Specifying this switch will turn this
                        behavior off.

  --timefmt=TIMEFMT     The format in which to display dates. %d - day, %b -
                        month, %m - month number, %Y - year. Default is: %b,
                        %Y

  --template=TEMPLATE   The template to control the filename and folder
                        structure of the saved files. Default is
                        "{author_sort}/{title}/{title} - {authors}" which will
                        save books into a per-author subfolder with filenames
                        containing title and author. Available controls are:
                        {author_sort, authors, id, isbn, languages,
                        last_modified, pubdate, publisher, rating, series,
                        series_index, tags, timestamp, title}

  --formats=FORMATS     Comma separated list of formats to save for each book.
                        By default all available formats are saved.

  --replace-whitespace  Replace whitespace with underscores.

  --to-lowercase        Convert paths to lowercase.


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

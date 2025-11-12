package calibredb
		
	// Usage: calibredb catalog /path/to/destination.(csv|epub|mobi|xml...) [options]
// 
// Export a catalog in format specified by path/to/destination extension.
// Options control how entries are displayed in the generated catalog output.
// Note that different catalog formats support different sets of options. To
// see the different options, specify the name of the output file and then the
// --help option.
// 
// 
// Whenever you pass arguments to calibredb that have spaces in them, enclose the arguments in quotation marks. For example: "/some path/with spaces"
// 
// Options:
//   -i IDS, --ids=IDS     Comma-separated list of database IDs to catalog.
//                         If declared, --search is ignored.
//                         Default: all
// 
//   -s SEARCH_TEXT, --search=SEARCH_TEXT
//                         Filter the results by the search query. For the format
//                         of the search query, please see the search-related
//                         documentation in the User Manual.
//                         Default: no filtering
// 
//   -v, --verbose         Show detailed output information. Useful for debugging
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
//   EPUB OPTIONS:
//     --catalog-title=CATALOG_TITLE
//                         Title of generated catalog used as title in metadata.
//                         Default: 'My Books'
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --cross-reference-authors
//                         Create cross-references in Authors section for books
//                         with multiple authors.
//                         Default: 'False'
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --debug-pipeline=DEBUG_PIPELINE
//                         Save the output from different stages of the
//                         conversion pipeline to the specified folder. Useful if
//                         you are unsure at which stage of the conversion
//                         process a bug is occurring.
//                         Default: 'none'
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --exclude-genre=EXCLUDE_GENRE
//                         Regex describing tags to exclude as genres.
//                         Default: '\[.+\]|^\+$' excludes bracketed tags, e.g.
//                         '[Project Gutenberg]', and '+', the default tag for
//                         read books.
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --exclusion-rules=EXCLUSION_RULES
//                         Specifies the rules used to exclude books from the
//                         generated catalog.
//                         The model for an exclusion rule is either
//                         ('<rule name>','Tags','<comma-separated list of
//                         tags>') or
//                         ('<rule name>','<custom column>','<pattern>').
//                         For example:
//                         (('Archived books','#status','Archived'),)
//                         will exclude a book with a value of 'Archived' in the
//                         custom column 'status'.
//                         When multiple rules are defined, all rules will be
//                         applied.
//                         Default:
//                         "(('Catalogs','Tags','Catalog'),)"
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --generate-authors  Include 'Authors' section in catalog.
//                         Default: 'False'
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --generate-descriptions
//                         Include 'Descriptions' section in catalog.
//                         Default: 'False'
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --generate-genres   Include 'Genres' section in catalog.
//                         Default: 'False'
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --generate-titles   Include 'Titles' section in catalog.
//                         Default: 'False'
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --generate-series   Include 'Series' section in catalog.
//                         Default: 'False'
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --generate-recently-added
//                         Include 'Recently Added' section in catalog.
//                         Default: 'False'
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --genre-source-field=GENRE_SOURCE_FIELD
//                         Source field for 'Genres' section.
//                         Default: 'Tags'
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --header-note-source-field=HEADER_NOTE_SOURCE_FIELD
//                         Custom field containing note text to insert in
//                         Description header.
//                         Default: ''
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --merge-comments-rule=MERGE_COMMENTS_RULE
//                         #<custom field>:[before|after]:[True|False]
//                         specifying:
//                          <custom field> Custom field containing notes to merge
//                         with comments
//                          [before|after] Placement of notes with respect to
//                         comments
//                          [True|False] - A horizontal rule is inserted between
//                         notes and comments
//                         Default: '::'
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --output-profile=OUTPUT_PROFILE
//                         Specifies the output profile. In some cases, an output
//                         profile is required to optimize the catalog for the
//                         device. For example, 'kindle' or 'kindle_dx' creates a
//                         structured Table of Contents with Sections and
//                         Articles.
//                         Default: 'none'
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --prefix-rules=PREFIX_RULES
//                         Specifies the rules used to include prefixes
//                         indicating read books, wishlist items and other user-
//                         specified prefixes.
//                         The model for a prefix rule is ('<rule name>','<source
//                         field>','<pattern>','<prefix>').
//                         When multiple rules are defined, the first matching
//                         rule will be used.
//                         Default:
//                         "(('Read books','tags','+','✓'),('Wishlist
//                         item','tags','Wishlist','×'))"
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --preset=PRESET     Use a named preset created with the GUI catalog
//                         builder.
//                         A preset specifies all settings for building a
//                         catalog.
//                         Default: 'none'
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --use-existing-cover
//                         Replace existing cover when generating the catalog.
//                         Default: 'False'
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
//     --thumb-width=THUMB_WIDTH
//                         Size hint (in inches) for book covers in catalog.
//                         Range: 1.0 - 2.0
//                         Default: '1.0'
//                         Applies to: AZW3, EPUB, MOBI output formats
// 
// 
// Created by Kovid Goyal <kovid@kovidgoyal.net>
// 
type CatalogOptions struct {
	GenerateAuthors string `json:"generate-authors,omitempty"`
	GenerateGenres string `json:"generate-genres,omitempty"`
	GenerateSeries string `json:"generate-series,omitempty"`
	GenerateTitles string `json:"generate-titles,omitempty"`
	Ids string `json:"ids,omitempty"`
	Preset string `json:"preset,omitempty"`
	Search string `json:"search,omitempty"`
	Verbose string `json:"verbose,omitempty"`
}

func (c *Calibre) CatalogHelp() string {
	return c.run("catalog", "-h")
}

func (c *Calibre) Catalog(opts CatalogOptions, args ...string) string {
	return "catalog"
}

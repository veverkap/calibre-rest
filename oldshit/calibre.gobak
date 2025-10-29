package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// CalibreWrapper wraps calibredb command-line operations
type CalibreWrapper struct {
	calibreDB    string
	library      string
	cdbWithLib   string
	logger       *log.Logger
	mutex        sync.Mutex
}

// Error types
type CalibreRuntimeError struct {
	Cmd        []string
	ReturnCode int
	Stdout     string
	Stderr     string
}

func (e *CalibreRuntimeError) Error() string {
	return fmt.Sprintf("Command %v failed with code %d: %s", e.Cmd, e.ReturnCode, e.Stderr)
}

type CalibreConcurrencyError struct {
	Cmd        []string
	ReturnCode int
}

func (e *CalibreConcurrencyError) Error() string {
	return fmt.Sprintf("Another calibre program is running (cmd: %v, code: %d)", e.Cmd, e.ReturnCode)
}

type ExistingItemError struct {
	Message string
}

func (e *ExistingItemError) Error() string {
	return e.Message
}

// Constants
var (
	addFlags = map[string]string{
		"authors":      "authors",
		"cover":        "cover",
		"identifiers":  "identifier",
		"isbn":         "isbn",
		"languages":    "languages",
		"series":       "series",
		"series_index": "series-index",
		"tags":         "tags",
		"title":        "title",
	}

	updateFlags = []string{
		"author_sort", "authors", "comments", "id", "identifiers",
		"languages", "pubdate", "publisher", "rating", "series",
		"series_index", "size", "tags", "timestamp", "title",
	}

	sortByKeys = []string{
		"author_sort", "authors", "comments", "cover", "formats",
		"id", "identifiers", "isbn", "languages", "last_modified",
		"pubdate", "publisher", "rating", "series", "series_index",
		"size", "tags", "template", "timestamp", "title", "uuid",
	}

	allowedFileExtensions = []string{
		".azw", ".azw3", ".azw4", ".cbz", ".cbr", ".cb7", ".cbc",
		".chm", ".djvu", ".docx", ".epub", ".fb2", ".fbz", ".html",
		".htmlz", ".lit", ".lrf", ".mobi", ".odt", ".pdf", ".prc",
		".pdb", ".pml", ".rb", ".rtf", ".snb", ".tcr", ".txt", ".txtz",
	}

	automergeValidValues = []string{"overwrite", "new_record", "ignore"}

	concurrencyErrRegex = regexp.MustCompile(`^Another calibre program.*is running.`)
	calibreVersionRegex = regexp.MustCompile(`calibre ([\d.]+)`)
	bookAddedRegex      = regexp.MustCompile(`^Added book ids: ([0-9, ]+)`)
	bookMergedRegex     = regexp.MustCompile(`^Merged book ids: ([0-9, ]+)`)
	bookIgnoredRegex    = regexp.MustCompile(`^The following books were not added as they already exist.*`)
)

// NewCalibreWrapper creates a new CalibreWrapper
func NewCalibreWrapper(calibreDB, library, username, password string, logger *log.Logger) *CalibreWrapper {
	if logger == nil {
		logger = log.New(os.Stdout, "[calibre] ", log.LstdFlags)
	}

	cdbWithLib := fmt.Sprintf("%s --with-library %s", calibreDB, library)
	if username != "" && password != "" {
		cdbWithLib += fmt.Sprintf(" --username %s --password %s", username, password)
	}

	return &CalibreWrapper{
		calibreDB:  calibreDB,
		library:    library,
		cdbWithLib: cdbWithLib,
		logger:     logger,
	}
}

// Check verifies that the executable and library exist
func (cw *CalibreWrapper) Check() error {
	// Check if calibredb executable exists
	if _, err := exec.LookPath(cw.calibreDB); err != nil {
		return fmt.Errorf("%s is not a valid executable: %w", cw.calibreDB, err)
	}

	// Check if metadata.db exists in library
	metadataPath := filepath.Join(cw.library, "metadata.db")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		return fmt.Errorf("failed to find Calibre database file %s", metadataPath)
	}

	return nil
}

// run executes a calibredb command
func (cw *CalibreWrapper) run(cmdStr string) (string, string, error) {
	cw.logger.Printf("Running: %s", cmdStr)
	
	cw.mutex.Lock()
	defer cw.mutex.Unlock()

	// Split command string into args
	args := strings.Fields(cmdStr)
	cmd := exec.Command(args[0], args[1:]...)

	stdout, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			stderr := string(exitError.Stderr)
			
			// Check for concurrency error
			if concurrencyErrRegex.MatchString(stderr) {
				return "", stderr, &CalibreConcurrencyError{
					Cmd:        args,
					ReturnCode: exitError.ExitCode(),
				}
			}
			
			return "", stderr, &CalibreRuntimeError{
				Cmd:        args,
				ReturnCode: exitError.ExitCode(),
				Stdout:     string(stdout),
				Stderr:     stderr,
			}
		}
		return "", "", err
	}

	return string(stdout), "", nil
}

// Version gets the calibredb version
func (cw *CalibreWrapper) Version() (string, error) {
	cmd := fmt.Sprintf("%s --version", cw.calibreDB)
	stdout, _, err := cw.run(cmd)
	if err != nil {
		return "", err
	}

	matches := calibreVersionRegex.FindStringSubmatch(stdout)
	if len(matches) > 1 {
		return matches[1], nil
	}

	cw.logger.Printf("failed to parse calibredb version from: %s", stdout)
	return "", fmt.Errorf("failed to parse calibredb version")
}

// GetBook gets a single book by ID
func (cw *CalibreWrapper) GetBook(id int) (*Book, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("%s list --for-machine --fields=all --search=id:%d --limit=1", cw.cdbWithLib, id)
	stdout, _, err := cw.run(cmd)
	if err != nil {
		return nil, err
	}

	var books []Book
	if err := json.Unmarshal([]byte(stdout), &books); err != nil {
		cw.logger.Printf("Error decoding JSON: %v\n\n%s", err, stdout)
		return nil, err
	}

	if len(books) == 1 {
		return &books[0], nil
	}

	return nil, nil
}

// GetBooks gets a list of books with optional sorting and searching
func (cw *CalibreWrapper) GetBooks(sort, search []string, all bool) ([]Book, error) {
	maxLimit := "all"
	if !all {
		maxLimit = "5000"
	}

	cmd := fmt.Sprintf("%s list --for-machine --fields=all --limit=%s", cw.cdbWithLib, maxLimit)
	cmd = cw.handleSort(cmd, sort)
	cmd = cw.handleSearch(cmd, search)

	stdout, _, err := cw.run(cmd)
	if err != nil {
		return nil, err
	}

	var books []Book
	if err := json.Unmarshal([]byte(stdout), &books); err != nil {
		return nil, err
	}

	return books, nil
}

// handleSort adds sort parameters to command
func (cw *CalibreWrapper) handleSort(cmd string, sort []string) string {
	if sort == nil || len(sort) == 0 {
		return cmd + " --ascending"
	}

	// Filter for supported sort keys
	safeSorts := make([]string, 0)
	unsafeSorts := make([]string, 0)

	for _, s := range sort {
		key := strings.TrimPrefix(s, "-")
		if contains(sortByKeys, key) {
			safeSorts = append(safeSorts, s)
		} else {
			unsafeSorts = append(unsafeSorts, s)
		}
	}

	if len(unsafeSorts) > 0 {
		cw.logger.Printf("The following sort keys are not supported and will be ignored: %q", strings.Join(unsafeSorts, ", "))
	}

	// Check if any sort key starts with "-" (descending)
	descending := false
	for _, s := range safeSorts {
		if strings.HasPrefix(s, "-") {
			descending = true
			break
		}
	}

	if !descending {
		cmd += " --ascending"
	}

	if len(safeSorts) > 0 {
		// Remove "-" prefix from sort keys
		cleanSorts := make([]string, len(safeSorts))
		for i, s := range safeSorts {
			cleanSorts[i] = strings.TrimPrefix(s, "-")
		}
		cmd += fmt.Sprintf(" --sort-by=%s", strings.Join(cleanSorts, ","))
	}

	return cmd
}

// handleSearch adds search parameters to command
func (cw *CalibreWrapper) handleSearch(cmd string, search []string) string {
	if search == nil || len(search) == 0 {
		return cmd
	}

	return cmd + fmt.Sprintf(` --search "%s"`, strings.Join(search, " "))
}

// validateID validates a book ID
func validateID(id int) error {
	if id <= 0 {
		return fmt.Errorf("book ID must be greater than 0")
	}
	return nil
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// AddMultiple adds multiple book files to the library
func (cw *CalibreWrapper) AddMultiple(filepaths []string, book Book, automerge string) ([]int, error) {
	if !contains(automergeValidValues, automerge) {
		return nil, fmt.Errorf("invalid automerge value: %s", automerge)
	}

	cmd := fmt.Sprintf("%s add", cw.cdbWithLib)
	
	// Add automerge flag
	if automerge != "ignore" {
		cmd += fmt.Sprintf(" --automerge=%s", automerge)
	}

	// Add book metadata flags
	cmd = cw.addBookFlags(cmd, book)

	// Add file paths
	for _, filepath := range filepaths {
		cmd += fmt.Sprintf(" %q", filepath)
	}

	stdout, _, err := cw.run(cmd)
	if err != nil {
		return nil, err
	}

	return cw.parseBookIDs(stdout)
}

// AddOneEmpty adds an empty book to the library
func (cw *CalibreWrapper) AddOneEmpty(book Book, automerge string) ([]int, error) {
	if !contains(automergeValidValues, automerge) {
		return nil, fmt.Errorf("invalid automerge value: %s", automerge)
	}

	cmd := fmt.Sprintf("%s add --empty", cw.cdbWithLib)
	
	// Add automerge flag
	if automerge != "ignore" {
		cmd += fmt.Sprintf(" --automerge=%s", automerge)
	}

	// Add book metadata flags
	cmd = cw.addBookFlags(cmd, book)

	stdout, _, err := cw.run(cmd)
	if err != nil {
		return nil, err
	}

	return cw.parseBookIDs(stdout)
}

// addBookFlags adds book metadata flags to command
func (cw *CalibreWrapper) addBookFlags(cmd string, book Book) string {
	if len(book.Authors) > 0 {
		cmd += fmt.Sprintf(" --authors %q", strings.Join(book.Authors, ","))
	}
	if book.Title != "" {
		cmd += fmt.Sprintf(" --title %q", book.Title)
	}
	if book.ISBN != "" {
		cmd += fmt.Sprintf(" --isbn %q", book.ISBN)
	}
	if len(book.Tags) > 0 {
		cmd += fmt.Sprintf(" --tags %q", strings.Join(book.Tags, ","))
	}
	if book.Series != "" {
		cmd += fmt.Sprintf(" --series %q", book.Series)
	}
	if book.SeriesIndex != 0 {
		cmd += fmt.Sprintf(" --series-index %f", book.SeriesIndex)
	}
	if book.Cover != "" {
		cmd += fmt.Sprintf(" --cover %q", book.Cover)
	}
	if len(book.Languages) > 0 {
		cmd += fmt.Sprintf(" --languages %q", strings.Join(book.Languages, ","))
	}
	
	// Handle identifiers
	for key, value := range book.Identifiers {
		cmd += fmt.Sprintf(" --identifier %s:%s", key, value)
	}

	return cmd
}

// parseBookIDs parses book IDs from command output
func (cw *CalibreWrapper) parseBookIDs(output string) ([]int, error) {
	// Check for added books
	if matches := bookAddedRegex.FindStringSubmatch(output); len(matches) > 1 {
		return parseIDList(matches[1])
	}
	
	// Check for merged books
	if matches := bookMergedRegex.FindStringSubmatch(output); len(matches) > 1 {
		return parseIDList(matches[1])
	}
	
	// Check for ignored books (existing)
	if bookIgnoredRegex.MatchString(output) {
		return nil, &ExistingItemError{Message: "Book already exists in library"}
	}

	return nil, fmt.Errorf("failed to parse book IDs from output: %s", output)
}

// parseIDList parses a comma-separated list of IDs
func parseIDList(idStr string) ([]int, error) {
	parts := strings.Split(strings.TrimSpace(idStr), ",")
	ids := make([]int, len(parts))
	
	for i, part := range parts {
		id, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("failed to parse ID %q: %w", part, err)
		}
		ids[i] = id
	}
	
	return ids, nil
}

// Remove removes books from the calibre database
func (cw *CalibreWrapper) Remove(ids []int, permanent bool) error {
	for _, id := range ids {
		if id <= 0 {
			return fmt.Errorf("invalid book ID: %d", id)
		}
	}

	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = strconv.Itoa(id)
	}

	cmd := fmt.Sprintf("%s remove %s", cw.cdbWithLib, strings.Join(idStrs, ","))
	if permanent {
		cmd += " --permanent"
	}

	_, _, err := cw.run(cmd)
	return err
}

// SetMetadata updates book metadata
func (cw *CalibreWrapper) SetMetadata(id int, book Book) error {
	if err := validateID(id); err != nil {
		return err
	}

	// Check if book exists
	existingBook, err := cw.GetBook(id)
	if err != nil {
		return err
	}
	if existingBook == nil {
		return fmt.Errorf("book with ID %d does not exist", id)
	}

	cmd := fmt.Sprintf("%s set_metadata %d", cw.cdbWithLib, id)
	cmd = cw.addUpdateFlags(cmd, book)

	_, _, err = cw.run(cmd)
	return err
}

// addUpdateFlags adds book metadata flags for set_metadata command
func (cw *CalibreWrapper) addUpdateFlags(cmd string, book Book) string {
	if len(book.Authors) > 0 {
		authors := strings.Join(book.Authors, " & ")
		cmd += fmt.Sprintf(" --field \"authors:%s\"", authors)
	}
	if book.Title != "" {
		cmd += fmt.Sprintf(" --field \"title:%s\"", book.Title)
	}
	if book.AuthorSort != "" {
		cmd += fmt.Sprintf(" --field \"author_sort:%s\"", book.AuthorSort)
	}
	if book.Comments != "" {
		cmd += fmt.Sprintf(" --field \"comments:%s\"", book.Comments)
	}
	if book.Publisher != "" {
		cmd += fmt.Sprintf(" --field \"publisher:%s\"", book.Publisher)
	}
	if book.PubDate != "" {
		cmd += fmt.Sprintf(" --field \"pubdate:%s\"", book.PubDate)
	}
	if book.Rating != 0 {
		cmd += fmt.Sprintf(" --field \"rating:%d\"", book.Rating)
	}
	if book.Series != "" {
		cmd += fmt.Sprintf(" --field \"series:%s\"", book.Series)
	}
	if book.SeriesIndex != 0 {
		cmd += fmt.Sprintf(" --field \"series_index:%f\"", book.SeriesIndex)
	}
	if len(book.Tags) > 0 {
		tags := strings.Join(book.Tags, ",")
		cmd += fmt.Sprintf(" --field \"tags:%s\"", tags)
	}
	if len(book.Languages) > 0 {
		languages := strings.Join(book.Languages, ",")
		cmd += fmt.Sprintf(" --field \"languages:%s\"", languages)
	}
	
	// Handle identifiers
	if len(book.Identifiers) > 0 {
		identifierPairs := make([]string, 0, len(book.Identifiers))
		for key, value := range book.Identifiers {
			identifierPairs = append(identifierPairs, fmt.Sprintf("%s:%s", key, value))
		}
		identifiers := strings.Join(identifierPairs, ",")
		cmd += fmt.Sprintf(" --field \"identifiers:%s\"", identifiers)
	}

	return cmd
}

// Export exports books from the calibre database to a directory
func (cw *CalibreWrapper) Export(ids []int, exportsDir string, formats []string) error {
	for _, id := range ids {
		if err := validateID(id); err != nil {
			return err
		}
	}

	cmd := fmt.Sprintf("%s export --dont-write-opf --dont-save-cover --single-dir", cw.cdbWithLib)
	
	if exportsDir != "" {
		cmd += fmt.Sprintf(" --to-dir=%s", exportsDir)
	}

	if formats != nil && len(formats) > 0 {
		cmd += fmt.Sprintf(" --formats %s", strings.Join(formats, ","))
	}

	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = strconv.Itoa(id)
	}
	cmd += fmt.Sprintf(" %s", strings.Join(idStrs, " "))

	_, stderr, err := cw.run(cmd)
	if err != nil {
		// Check for "No book with id" error
		if strings.Contains(stderr, "No book with id") {
			return fmt.Errorf("one or more books not found: %s", stderr)
		}
		return err
	}

	return nil
}
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Server holds the HTTP server and dependencies
type Server struct {
	calibre *CalibreWrapper
	logger  *log.Logger
	version string
}

// NewServer creates a new HTTP server
func NewServer(calibre *CalibreWrapper, logger *log.Logger, version string) *Server {
	return &Server{
		calibre: calibre,
		logger:  logger,
		version: version,
	}
}

// SetupRoutes sets up HTTP routes
func (s *Server) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Health endpoint
	r.HandleFunc("/health", s.handleHealth).Methods("GET")

	// Book endpoints
	r.HandleFunc("/books/{id:[0-9]+}", s.handleGetBook).Methods("GET")
	r.HandleFunc("/books", s.handleGetBooks).Methods("GET")
	r.HandleFunc("/books", s.handleAddBook).Methods("POST")
	r.HandleFunc("/books/empty", s.handleAddEmptyBook).Methods("POST")
	r.HandleFunc("/books/{id:[0-9]+}", s.handleUpdateBook).Methods("PUT")
	r.HandleFunc("/books/{id:[0-9]+}", s.handleDeleteBook).Methods("DELETE")
	r.HandleFunc("/books", s.handleDeleteBooks).Methods("DELETE")

	// Export endpoints
	r.HandleFunc("/export/{id:[0-9]+}", s.handleExportBook).Methods("GET")
	r.HandleFunc("/export", s.handleExportBooks).Methods("GET")

	return r
}

// handleHealth returns version information
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	calibreVersion, err := s.calibre.Version()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get calibre version: %v", err))
		return
	}

	response := map[string]string{
		"calibre_version":      calibreVersion,
		"calibre_rest_version": s.version,
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleGetBook returns a single book by ID
func (s *Server) handleGetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	book, err := s.calibre.GetBook(id)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if book == nil {
		s.writeError(w, http.StatusNotFound, fmt.Sprintf("book %d does not exist", id))
		return
	}

	response := map[string]*Book{"books": book}
	s.writeJSON(w, http.StatusOK, response)
}

// handleGetBooks returns a paginated list of books
func (s *Server) handleGetBooks(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	startStr := r.URL.Query().Get("start")
	if startStr == "" {
		startStr = "1"
	}
	start, err := strconv.Atoi(startStr)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid start parameter")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "20"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	sort := r.URL.Query()["sort"]
	search := r.URL.Query()["search"]

	books, err := s.calibre.GetBooks(sort, search, false)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(books) == 0 {
		response := map[string][]Book{"books": {}}
		w.WriteHeader(http.StatusNoContent)
		s.writeJSON(w, http.StatusNoContent, response)
		return
	}

	paginated, err := NewPaginatedResults(books, start, limit, sort, search)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, paginated.ToDict())
}

// handleAddBook adds a book from uploaded file(s) with optional metadata
func (s *Server) handleAddBook(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		s.writeError(w, http.StatusUnsupportedMediaType, "Only multipart/form-data allowed")
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	// Check files
	filepaths, err := s.checkFiles(r)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer s.cleanupFiles(filepaths)

	// Parse optional JSON data
	var book Book
	automerge := "ignore"

	if dataStr := r.FormValue("data"); dataStr != "" {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
			s.writeError(w, http.StatusBadRequest, "Invalid JSON data")
			return
		}

		// Extract automerge if present
		if am, ok := data["automerge"].(string); ok {
			automerge = am
			delete(data, "automerge")
		}

		// Convert remaining data to Book struct
		if len(data) > 0 {
			dataBytes, _ := json.Marshal(data)
			if err := json.Unmarshal(dataBytes, &book); err != nil {
				s.writeError(w, http.StatusBadRequest, "Invalid book data")
				return
			}
		}
	}

	ids, err := s.calibre.AddMultiple(filepaths, book, automerge)
	if err != nil {
		if _, ok := err.(*ExistingItemError); ok {
			s.writeError(w, http.StatusConflict, err.Error())
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response := map[string][]int{"id": ids}
	s.writeJSON(w, http.StatusCreated, response)
}

// handleAddEmptyBook adds an empty book with metadata
func (s *Server) handleAddEmptyBook(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		s.writeError(w, http.StatusUnsupportedMediaType, "Only application/json allowed")
		return
	}

	var data map[string]interface{}
	
	// Allow empty request body
	if r.ContentLength > 0 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.writeError(w, http.StatusBadRequest, "Failed to read request body")
			return
		}

		if err := json.Unmarshal(body, &data); err != nil {
			s.writeError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
	}

	// Extract automerge if present
	automerge := "ignore"
	if am, ok := data["automerge"].(string); ok {
		automerge = am
		delete(data, "automerge")
	}

	// Convert to Book struct
	var book Book
	if len(data) > 0 {
		dataBytes, _ := json.Marshal(data)
		if err := json.Unmarshal(dataBytes, &book); err != nil {
			s.writeError(w, http.StatusBadRequest, "Invalid book data")
			return
		}
	}

	ids, err := s.calibre.AddOneEmpty(book, automerge)
	if err != nil {
		if _, ok := err.(*ExistingItemError); ok {
			s.writeError(w, http.StatusConflict, err.Error())
		} else {
			s.writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response := map[string][]int{"id": ids}
	s.writeJSON(w, http.StatusCreated, response)
}

// handleUpdateBook updates an existing book's metadata
func (s *Server) handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		s.writeError(w, http.StatusUnsupportedMediaType, "Only application/json allowed")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	if len(body) == 0 {
		s.writeError(w, http.StatusBadRequest, "No data provided")
		return
	}

	var book Book
	if err := json.Unmarshal(body, &book); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// For now, return not implemented as the calibre wrapper doesn't have SetMetadata
	s.writeError(w, http.StatusNotImplemented, "Update book not yet implemented")
}

// handleDeleteBook deletes a book by ID
func (s *Server) handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	// For now, return not implemented as the calibre wrapper doesn't have Remove
	s.writeError(w, http.StatusNotImplemented, "Delete book not yet implemented")
}

// handleDeleteBooks deletes multiple books by IDs
func (s *Server) handleDeleteBooks(w http.ResponseWriter, r *http.Request) {
	// For now, return not implemented
	s.writeError(w, http.StatusNotImplemented, "Delete books not yet implemented")
}

// handleExportBook exports a single book
func (s *Server) handleExportBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	// For now, return not implemented as the calibre wrapper doesn't have Export
	s.writeError(w, http.StatusNotImplemented, "Export book not yet implemented")
}

// handleExportBooks exports multiple books
func (s *Server) handleExportBooks(w http.ResponseWriter, r *http.Request) {
	// For now, return not implemented
	s.writeError(w, http.StatusNotImplemented, "Export books not yet implemented")
}

// checkFiles validates and saves uploaded files
func (s *Server) checkFiles(r *http.Request) ([]string, error) {
	if r.MultipartForm == nil || len(r.MultipartForm.File) == 0 {
		return nil, fmt.Errorf("no file(s) provided")
	}

	var filepaths []string
	tempDir := os.TempDir()

	// Process all uploaded files
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			if fileHeader.Filename == "" {
				return nil, fmt.Errorf("invalid file or filename")
			}

			if !s.allowedFile(fileHeader.Filename) {
				return nil, fmt.Errorf("invalid filename (%s)", fileHeader.Filename)
			}

			// Open uploaded file
			file, err := fileHeader.Open()
			if err != nil {
				return nil, err
			}
			defer file.Close()

			// Create temporary file
			tempFilePath := filepath.Join(tempDir, secureFilename(fileHeader.Filename))
			tempFile, err := os.Create(tempFilePath)
			if err != nil {
				return nil, err
			}
			defer tempFile.Close()

			// Copy file content
			if _, err := io.Copy(tempFile, file); err != nil {
				return nil, err
			}

			filepaths = append(filepaths, tempFilePath)
		}
	}

	return filepaths, nil
}

// allowedFile checks if a filename has an allowed extension
func (s *Server) allowedFile(filename string) bool {
	if strings.HasPrefix(filename, "-") {
		return false
	}

	filename = strings.ToLower(filename)
	for _, ext := range allowedFileExtensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

// secureFilename sanitizes a filename
func secureFilename(filename string) string {
	// Basic implementation - in production, use a more robust solution
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, "..", "_")
	return filename
}

// cleanupFiles removes temporary files
func (s *Server) cleanupFiles(filepaths []string) {
	for _, path := range filepaths {
		if err := os.Remove(path); err != nil {
			s.logger.Printf("Failed to remove temporary file %s: %v", path, err)
		}
	}
}

// writeJSON writes a JSON response
func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Printf("Failed to encode JSON response: %v", err)
	}
}

// writeError writes an error response
func (s *Server) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Printf("Failed to encode error response: %v", err)
	}
}
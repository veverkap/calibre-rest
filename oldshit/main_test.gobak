package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	// Create a mock calibre wrapper
	config := NewTestConfig()
	config.SetCalibreDB("/bin/echo") // Use echo as a mock calibredb for testing
	config.SetLibrary("/tmp")
	
	// Create a temporary metadata.db file for testing
	metadataPath := "/tmp/metadata.db"
	file, err := os.Create(metadataPath)
	if err != nil {
		t.Fatalf("Failed to create test metadata.db: %v", err)
	}
	file.Close()
	defer os.Remove(metadataPath)

	calibre := NewCalibreWrapper(config.CalibreDB, config.Library, "", "", nil)
	server := NewServer(calibre, nil, "test")

	router := server.SetupRoutes()

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if _, ok := response["calibre_rest_version"]; !ok {
		t.Error("Response missing calibre_rest_version")
	}
}

func TestGetBookNotFound(t *testing.T) {
	config := NewTestConfig()
	config.SetCalibreDB("/bin/echo")
	config.SetLibrary("/tmp")
	
	// Create a temporary metadata.db file for testing
	metadataPath := "/tmp/metadata.db"
	file, err := os.Create(metadataPath)
	if err != nil {
		t.Fatalf("Failed to create test metadata.db: %v", err)
	}
	file.Close()
	defer os.Remove(metadataPath)

	calibre := NewCalibreWrapper(config.CalibreDB, config.Library, "", "", nil)
	server := NewServer(calibre, nil, "test")

	router := server.SetupRoutes()

	req, err := http.NewRequest("GET", "/books/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Should return error since we're using echo which won't return valid JSON
	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler should not return OK for non-existent book")
	}
}

func TestNotImplementedEndpoints(t *testing.T) {
	config := NewTestConfig()
	config.SetCalibreDB("/bin/echo")
	config.SetLibrary("/tmp")
	
	// Create a temporary metadata.db file for testing
	metadataPath := "/tmp/metadata.db"
	file, err := os.Create(metadataPath)
	if err != nil {
		t.Fatalf("Failed to create test metadata.db: %v", err)
	}
	file.Close()
	defer os.Remove(metadataPath)

	calibre := NewCalibreWrapper(config.CalibreDB, config.Library, "", "", nil)
	server := NewServer(calibre, nil, "test")

	router := server.SetupRoutes()

	// Test endpoints that should return errors (but not 501 anymore since they're implemented)
	testCases := []struct {
		method string
		path   string
		expectStatus int
	}{
		{"PUT", "/books/1", http.StatusInternalServerError}, // Will fail due to echo not returning valid book
		{"DELETE", "/books/1", http.StatusInternalServerError}, // Will fail during verification
		{"GET", "/export/1", http.StatusInternalServerError}, // Will fail during export
	}

	for _, tc := range testCases {
		req, err := http.NewRequest(tc.method, tc.path, strings.NewReader("{}"))
		if err != nil {
			t.Fatal(err)
		}

		if tc.method == "PUT" {
			req.Header.Set("Content-Type", "application/json")
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != tc.expectStatus {
			t.Errorf("%s %s returned status code: got %v want %v", tc.method, tc.path, status, tc.expectStatus)
		}
	}
}
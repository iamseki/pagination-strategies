package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func TestBooksOffsetHandler(t *testing.T) {
	// Connect to PostgreSQL database using sqlx
	db, err := sqlx.Connect("pgx", "postgres://test:test@localhost:5432/library?sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	// Start a server in-memory
	router := mux.NewRouter()
	router.HandleFunc("/books/offset", func(w http.ResponseWriter, r *http.Request) { BooksOffsetHandler(db, w, r) })
	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("First page (limit 10)", func(t *testing.T) {
		url := fmt.Sprintf("%s/books/offset?limit=10", server.URL)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal("error while retrieving the first page", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", res.StatusCode)
		}

		var firstPage PagedOffsetResponse
		if err := json.NewDecoder(res.Body).Decode(&firstPage); err != nil {
			t.Fatal("error while decode response body", err)
		}

		// Assert response length (check if 10 books are returned)
		if len(firstPage.Books) != 10 {
			t.Errorf("Expected 10 books, got %d", len(firstPage.Books))
		}
	})

	t.Run("Second page (limit 10)", func(t *testing.T) {
		// Test case 2: Retrieve second page (limit 10, offset 10)
		url := fmt.Sprintf("%s/books/offset?limit=10&offset=10", server.URL)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", res.StatusCode)
		}

		// Decode resonse body
		var secondPage PagedOffsetResponse
		if err := json.NewDecoder(res.Body).Decode(&secondPage); err != nil {
			t.Fatal(err)
		}

		// Assert response length (check if 10 books are returned)
		if len(secondPage.Books) != 10 {
			t.Errorf("Expected 10 books, got %d", len(secondPage.Books))
		}

	})

}

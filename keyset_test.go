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

func TestBooksKeysetHandler(t *testing.T) {
	// Connect to PostgreSQL database using sqlx
	db, err := sqlx.Connect("pgx", "postgres://test:test@localhost:5432/library?sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	// Start a server in-memory
	router := mux.NewRouter()
	router.HandleFunc("/books/keyset", func(w http.ResponseWriter, r *http.Request) { BooksKeysetHandler(db, w, r) })
	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("Retrieve first two pages with limit 10", func(t *testing.T) {
		url := fmt.Sprintf("%s/books/keyset?limit=10", server.URL)
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

		var firstPage PagedKeysetResponse
		if err := json.NewDecoder(res.Body).Decode(&firstPage); err != nil {
			t.Fatal("error while decode response body", err)
		}

		// Assert response length (check if 10 books are returned)
		if len(firstPage.Books) != 10 {
			t.Errorf("Expected 10 books, got %d", len(firstPage.Books))
		}

		if firstPage.Books[len(firstPage.Books)-1].ID != 10 {
			t.Errorf("Expected last book from first page to be ID = 10, got %d", firstPage.Books[len(firstPage.Books)-1].ID)
		}

		// Second Page
		url = fmt.Sprintf("%s/books/keyset?limit=10&nextToken=%s", server.URL, firstPage.NextToken)
		req, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal("error while retrieving the first page", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", res.StatusCode)
		}

		var secondPage PagedKeysetResponse
		if err = json.NewDecoder(res.Body).Decode(&secondPage); err != nil {
			t.Fatal("error while decode response body", err)
		}

		if len(secondPage.Books) != 10 {
			t.Errorf("Expected 10 books, got %d", len(secondPage.Books))
		}

		if secondPage.Books[0].ID != 11 {
			t.Errorf("Expected first book from second page to be ID = 11, got %d", secondPage.Books[0].ID)
		}

		if secondPage.Books[len(secondPage.Books)-1].ID != 20 {
			t.Errorf("Expected last book from second page to be ID = 20, got %d", secondPage.Books[len(secondPage.Books)-1].ID)
		}
	})
}

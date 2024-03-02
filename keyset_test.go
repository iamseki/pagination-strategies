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

func TestDecodeToken(t *testing.T) {
	token := encodeToken(17, Backward)
	id, direction := decodeToken(token)
	if id != 17 {
		t.Errorf("expecet id to be 17, got: %v", id)
	}

	if direction != Backward {
		t.Errorf("expect direction to be Backward, got: %v", direction)
	}
}

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

	t.Run("Move three pages forward and then back one (nextPageToken + lastPageToken)", func(t *testing.T) {
		{
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

			var firstPage PagedKeysetResponse
			if err := json.NewDecoder(res.Body).Decode(&firstPage); err != nil {
				t.Fatal("error while decode response body", err)
			}

			// Second Page
			url = fmt.Sprintf("%s/books/keyset?limit=10&pageToken=%s", server.URL, firstPage.NextToken)
			req, err = http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err = http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal("error while retrieving the first page", err)
			}
			defer res.Body.Close()

			var secondPage PagedKeysetResponse
			if err = json.NewDecoder(res.Body).Decode(&secondPage); err != nil {
				t.Fatal("error while decode response body", err)
			}

			// Third page
			url = fmt.Sprintf("%s/books/keyset?limit=10&pageToken=%s", server.URL, secondPage.NextToken)
			req, err = http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err = http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal("error while retrieving the first page", err)
			}
			defer res.Body.Close()

			var thirdPage PagedKeysetResponse
			if err = json.NewDecoder(res.Body).Decode(&thirdPage); err != nil {
				t.Fatal("error while decode response body", err)
			}

			// Backward
			url = fmt.Sprintf("%s/books/keyset?limit=10&pageToken=%s", server.URL, thirdPage.PreviousToken)
			req, err = http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err = http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal("error while retrieving the first page", err)
			}
			defer res.Body.Close()

			var previousPage PagedKeysetResponse
			if err = json.NewDecoder(res.Body).Decode(&previousPage); err != nil {
				t.Fatal("error while decode response body", err)
			}

			if len(previousPage.Books) != 10 {
				t.Errorf("Expected 10 books, got %d", len(previousPage.Books))
			}

			if previousPage.Books[0].ID != 20 {
				t.Errorf("Expected first book from previous page to be ID = 20, got %d", previousPage.Books[0].ID)
			}

			if previousPage.Books[len(previousPage.Books)-1].ID != 11 {
				t.Errorf("Expected last book from previous page to be ID = 11, got %d", previousPage.Books[len(previousPage.Books)-1].ID)
			}
		}
	})

	t.Run("Retrieve first two pages with limit 10 scan forward (nextPageToken only)", func(t *testing.T) {
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
		url = fmt.Sprintf("%s/books/keyset?limit=10&pageToken=%s", server.URL, firstPage.NextToken)
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

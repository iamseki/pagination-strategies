package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
)

func BooksKeysetHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	// for now the token it's the id itself
	nextKeysetToken := r.URL.Query().Get("nextToken")
	book_id := 0
	if nextKeysetToken != "" {
		decoded, err := base64.StdEncoding.DecodeString(nextKeysetToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Fail to decode nextKeysetToken")
			return
		}
		book_id, err = strconv.Atoi(string(decoded))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid nextKeysetToken")
			return
		}
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid limit parameter")
		return
	}

	books, err := retriveBooksKeyset(db, limit, book_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error on retriveBooksKeyset: %v", err)
		return
	}

	// Encode books to JSON and write to response
	nextToken := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(books[len(books)-1].ID)))
	res := PagedKeysetResponse{
		Books:     books,
		NextToken: nextToken,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding books to JSON: %v", err)
		return
	}
}

// HashKeyset securely hashes a string using SHA-256
func hashKeyset(keyset string) string {
	hash := sha256.Sum256([]byte(keyset))
	return string(hash[:])
}

func retriveBooksKeyset(db *sqlx.DB, limit int, id int) ([]Book, error) {
	query := `SELECT * FROM books WHERE id > $1 ORDER BY id ASC LIMIT $2`

	rows, err := db.Query(query, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

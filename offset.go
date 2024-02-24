package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
)

func BooksOffsetHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	// Extract offset and limit parameters from request query (adjust as needed)
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid offset parameter")
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid limit parameter")
		return
	}

	// Connect to database and retrieve books with offset and limit
	// Replace with your actual database connection and query logic
	// Remember to close the connection after usage

	books, err := retrieveBooksOffset(db, offset, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error retrieving books: %v", err)
		return
	}

	// Encode books to JSON and write to response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding books to JSON: %v", err)
		return
	}
}

func retrieveBooksOffset(db *sqlx.DB, offset, limit int) ([]Book, error) {
	// Replace with your actual query using sqlx syntax
	query := `SELECT * FROM books LIMIT $1 OFFSET $2`
	rows, err := db.Query(query, limit, offset)
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

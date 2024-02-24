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
	offsetStr := r.URL.Query().Get("offset")
	var offset int
	if offsetStr == "" {
		offset = 0
	} else {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid offset parameter")
			return
		}
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid limit parameter")
		return
	}

	total, err := countBooksOffset(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error count books: %v", err)
		return
	}

	books, err := retrieveBooksOffset(db, offset, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error retrieving books: %v", err)
		return
	}

	next := ""
	if len(books) == limit && total > offset+limit {
		nextURL := fmt.Sprintf("/books/offset?offset=%d&limit=%d", offset+limit, limit)
		next = nextURL
	}

	prev := ""
	if offset > 0 {
		prevURL := fmt.Sprintf("/books/offset?offset=%d&limit=%d", max(offset-limit, 0), limit)
		prev = prevURL
	}

	res := PagedOffsetResponse{
		Books: books,
		Next:  next,
		Prev:  prev,
		Count: total,
	}

	// Encode books to JSON and write to response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding books to JSON: %v", err)
		return
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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

func countBooksOffset(db *sqlx.DB) (int, error) {
	total := 0
	err := db.Get(&total, "SELECT count(*) FROM books")
	if err != nil {
		return total, err
	}
	return total, nil
}

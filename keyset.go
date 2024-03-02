package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Direction string

const (
	Forward  Direction = "forward"
	Backward Direction = "backward"
)

// Decodes a token string using base64 extracting its id and Direction
func decodeToken(token string) (int, Direction) {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		log.Fatalf("error decodeToken: %v", err)
	}

	values := strings.Split(string(decoded), ",")
	idStr, direction := values[0], values[1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatalf("error decodeToken convert id to int: %v", err)
	}

	return id, Direction(direction)
}

func encodeToken(id int, d Direction) string {
	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d,%s", id, d)))
	return token
}

func BooksKeysetHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid limit parameter")
		return
	}
	// default to first page with direction forward
	last_page_id := 0
	direction := Forward
	pageToken := r.URL.Query().Get("pageToken")

	if pageToken != "" {
		last_page_id, direction = decodeToken(pageToken)
	}

	var nextToken, previousToken string
	var books []Book

	switch direction {
	case Backward:
		books, err = retriveBooksKeysetBackward(db, limit+1, last_page_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error on retriveBooksKeyset: %v", err)
			return
		}

		hasMorePages := len(books) > limit
		if hasMorePages {
			// create a new slice with all elements except the LAST ONE
			books = books[:len(books)-1]
			previousToken = encodeToken(books[len(books)-1].ID, Backward)
		}
		nextToken = encodeToken(books[0].ID, Forward)

	case Forward:
		books, err = retriveBooksKeysetForward(db, limit+1, last_page_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error on retriveBooksKeyset: %v", err)
			return
		}

		hasMorePages := len(books) > limit
		if hasMorePages {
			// create a new slice with all elements except the LAST ONE
			books = books[:len(books)-1]
			nextToken = encodeToken(books[len(books)-1].ID, Forward)
		}
		// in first page scenario previousToken must be nil
		if pageToken != "" {
			previousToken = encodeToken(books[0].ID, Backward)
		}
	}

	res := PagedKeysetResponse{
		Books:         books,
		NextToken:     nextToken,
		PreviousToken: previousToken,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding books to JSON: %v", err)
		return
	}
}

func retriveBooksKeysetForward(db *sqlx.DB, limit int, id int) ([]Book, error) {
	query := `SELECT * FROM books WHERE id > $1 ORDER BY id ASC LIMIT $2`

	rows, err := db.Query(query, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []Book{}
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

func retriveBooksKeysetBackward(db *sqlx.DB, limit int, id int) ([]Book, error) {
	query := `SELECT * FROM books WHERE id < $1 ORDER BY id DESC LIMIT $2`

	rows, err := db.Query(query, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []Book{}
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

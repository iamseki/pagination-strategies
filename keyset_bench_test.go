package main

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
)

func BenchmarkKeysetForwardScan(b *testing.B) {
	// Connect to PostgreSQL database using sqlx
	db, err := sqlx.Connect("pgx", "postgres://test:test@localhost:5432/library?sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	limit := 10
	id := 0

	books := []Book{}
	books, err = retriveBooksKeysetForward(db, limit, id)
	booksToQuery := 10000
	booksScanned := len(books)
	for booksScanned < booksToQuery {
		nextId := books[len(books)-1].ID
		books, _ = retriveBooksKeysetForward(db, limit, nextId)
		booksScanned += len(books)
	}
}

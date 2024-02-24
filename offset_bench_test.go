package main

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
)

func BenchmarkOffsetScan(b *testing.B) {
	// Connect to PostgreSQL database using sqlx
	db, err := sqlx.Connect("pgx", "postgres://test:test@localhost:5432/library?sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	limit := 10
	offset := 0

	books := []Book{}
	books, err = retrieveBooksOffset(db, offset, limit)
	booksToQuery := 10000
	booksScanned := len(books)
	for booksScanned < booksToQuery {
		offset += limit
		books, _ = retrieveBooksOffset(db, offset, limit)
		booksScanned += len(books)
	}
}

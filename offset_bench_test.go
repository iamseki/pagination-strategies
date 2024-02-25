package main

import (
	"fmt"
	"log"
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

	for i := 0; i < b.N; i++ {
		limit := 10
		offset := 0
		scan10000BooksWithOffset(db, offset, limit)
	}

}

func scan10000BooksWithOffset(db *sqlx.DB, offset, limit int) {
	booksToQuery := 10000
	books := []Book{}
	var err error

	books, err = retrieveBooksOffset(db, offset, limit)
	if err != nil {
		log.Fatalf("error running scan10000BooksWithOffset: %v", err)
	}

	booksScanned := len(books)
	for booksScanned < booksToQuery {
		offset += limit
		books, _ = retrieveBooksOffset(db, offset, limit)
		booksScanned += len(books)
	}
}

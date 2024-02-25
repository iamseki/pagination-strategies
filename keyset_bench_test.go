package main

import (
	"fmt"
	"log"
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

	for i := 0; i < b.N; i++ {
		limit := 10
		id := 0
		scan10000BooksWithKeyset(db, id, limit)
	}
}

func scan10000BooksWithKeyset(db *sqlx.DB, id, limit int) {
	booksToQuery := 10000
	books := []Book{}
	var err error
	books, err = retriveBooksKeysetForward(db, limit, id)
	if err != nil {
		log.Fatalf("error on scan10000BooksWithKeyset: %v", err)
	}

	booksScanned := len(books)
	for booksScanned < booksToQuery {
		nextId := books[len(books)-1].ID
		books, _ = retriveBooksKeysetForward(db, limit, nextId)
		booksScanned += len(books)
	}
}

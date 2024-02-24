package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib" // Setup Postgres driver
	"github.com/jmoiron/sqlx"
)

// Book struct represents a book in your database
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Genre  string `json:"genre"`
}

type PagedResponse struct {
	Books []Book `json:"books"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
	Count int    `json:"total"`
}

func main() {
	// Connect to PostgreSQL database using sqlx
	db, err := sqlx.Connect("pgx", "postgres://test:test@localhost:5432/library?sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	// Initialize router and register handler
	router := mux.NewRouter()
	router.HandleFunc("/books/offset", func(w http.ResponseWriter, r *http.Request) { BooksOffsetHandler(db, w, r) })
	// Start server on port 8080 (adjust as needed)
	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

package handlers

import (
	"book-managment/db"
	"book-managment/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GET /books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	limitParam := r.URL.Query().Get("limit")
	var rows *sql.Rows
	var err error

	if limitParam != "" {
		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
		rows, err = db.GetDb().Query("SELECT id, title, isbn, author, year FROM books LIMIT ?", limit)
	} else {
		rows, err = db.GetDb().Query("SELECT id, title, isbn, author, year FROM books")
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.ISBN, &book.Author, &book.Year); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}
	json.NewEncoder(w).Encode(books)
}

// GET /books/{id}
func GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var book models.Book
	err := db.GetDb().QueryRow("SELECT id, title, isbn, author, year FROM books WHERE id = ?", id).Scan(&book.ID, &book.Title, &book.ISBN, &book.Author, &book.Year)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(book)
}

// POST /books
func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.GetDb().Exec("INSERT INTO books (title, isbn, author, year) VALUES (?, ?, ?, ?)", book.Title, book.ISBN, book.Author, book.Year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	book.ID = uint(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// PUT /books/{id}
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.GetDb().Exec("UPDATE books SET title = ?, isbn = ?, author = ?, year = ? WHERE id = ?", book.Title, book.ISBN, book.Author, book.Year, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	book.ID = uint(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// DELETE /books/{id}
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	_, err := db.GetDb().Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

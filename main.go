package main

import (
	"book-managment/db"
	"book-managment/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Инициализация на базата данни
	db.Init()

	// Инициализация на маршрутизатора
	router := mux.NewRouter()
	router.HandleFunc("/books", handlers.GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", handlers.GetBook).Methods("GET")
	router.HandleFunc("/books", handlers.CreateBook).Methods("POST")
	router.HandleFunc("/books/{id}", handlers.UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", handlers.DeleteBook).Methods("DELETE")

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

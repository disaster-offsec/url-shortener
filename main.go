package main

import (
	"log"
	"net/http"
)

func main() {
	storage := NewStorage()

	// Пока что без API
	http.HandleFunc("POST /shorten", shortenHandler(storage))
	http.HandleFunc("GET /{short}", redirectHandler(storage))

	// Запускаем сервер на 8080 порт
	log.Println("Server starting on :8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func shortenHandler(storage *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Заглушка
		panic("shortenHandler not implemented yet")
	}
}

func redirectHandler(storage *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Заглушка
		panic("redirectHandler not implemented yet")
	}
}

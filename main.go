package main

import (
	"log"
	"net/http"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
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

// на первое время сгодится
func generateShortCode() string {
	b := make([]byte, 4)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:6]
}


// Проверка метода необязательна
func shortenHandler(storage *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			URL string `json:"url"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		code := generateShortCode()

		for storage.ExistCode(code) {
			code = generateShortCode()
		}

		storage.Save(code, req.URL)

		shortURL := "http://localhost:8080/" + code
		response := map[string]string{"short_url": shortURL}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func redirectHandler(storage *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("short")
		if code == "" {
			http.NotFound(w, r)
			return
		}
		
		// амортизированное O(1)
		originalURL, ok := storage.GetByCode(code)

		if !ok {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}

package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"context"
	"time"
)

func main() {
	storage := NewStorage()

	// Пока что без API
	http.HandleFunc("POST /shorten", shortenHandler(storage))
	http.HandleFunc("GET /{short}", redirectHandler(storage))

	// Запускаем сервер на 8080 порт
	log.Println("Server starting on :8080...")

	server := &http.Server{Addr: ":8080", Handler: nil}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeOut(context.Background(), time.Second * 5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exited properly")
}

func isValidURL(raw string) bool {
	parsed, err := url.Parse(raw)
	return err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https")
}

// Создаем криптографически стойкий 6-значный URL-безопасный-код
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

		if !isValidURL(req.URL) {
			http.Error(w, "Invalid URL scheme (must be http or https)", http.StatusBadRequest)
			return
		}	

		if existingCode, ok := storage.GetByOriginal(req.URL); ok {
			shortURL := "http://localhost:8080/" + existingCode
			response := map[string]string{"short_url": shortURL}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
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

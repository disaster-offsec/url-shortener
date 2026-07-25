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

	port := os.Getenv("PORT")
	if port == "" {port = "8080"}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {baseURL = "http://localhost:" + port}
	
	http.HandleFunc("POST /shorten", loggingMiddleware(shortenHandler(baseURL, storage)))
	http.HandleFunc("GET /{short}", loggingMiddleware(redirectHandler(storage)))


	log.Printf("Server starting on: %s", baseURL)

	server := &http.Server{Addr: ":" + port, Handler: nil}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exited properly")
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	}
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
func shortenHandler(baseURL string, storage *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			URL string `json:"url"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Invalid JSON: %v", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if !isValidURL(req.URL) {
			log.Printf("Invalid URL scheme (must be http or https)")
			http.Error(w, "Invalid URL scheme (must be http or https)", http.StatusBadRequest)
			return
		}	

		if existingCode, ok := storage.GetByOriginal(req.URL); ok {
			shortURL := baseURL + "/" + existingCode
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

		shortURL := baseURL + "/" + code
		response := map[string]string{"short_url": shortURL}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func redirectHandler(storage *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("short")
		if code == "" {
			log.Printf("Short code not found %s", code) // да там "", но тем не менее
			http.NotFound(w, r)
			return
		}

		// амортизированное O(1)
		originalURL, ok := storage.GetByCode(code)

		if !ok {
			log.Printf("Short code not found %s", code)
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}

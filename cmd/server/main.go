package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/handlers"
	"url-shortener/internal/storage"
)

func main() {
	cfg := config.Load()

	store := storage.NewStorage()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", handlers.LoggingMiddleware(handlers.ShortenHandler(cfg.BaseURL, store)))
	mux.HandleFunc("GET /{short}", handlers.LoggingMiddleware(handlers.RedirectHandler(store)))



	server := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		log.Printf("Server starting on: %s", cfg.BaseURL)
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

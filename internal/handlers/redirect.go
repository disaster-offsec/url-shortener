package handlers

import (
	"log"
	"net/http"

	"url-shortener/internal/storage"
)

func RedirectHandler(store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("short")
		if code == "" {
			log.Printf("Short code not found %s", code) // да там "", но тем не менее
			http.NotFound(w, r)
			return
		}

		// амортизированное O(1)
		originalURL, ok := store.GetByCode(code)

		if !ok {
			log.Printf("Short code not found %s", code)
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}

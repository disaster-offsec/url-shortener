package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"url-shortener/internal/storage"
	"url-shortener/internal/utils"
)

func ShortenHandler(baseURL string, store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			URL string `json:"url"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Invalid JSON: %v", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if !utils.IsValidURL(req.URL) {
			log.Printf("Invalid URL scheme (must be http or https)")
			http.Error(w, "Invalid URL scheme (must be http or https)", http.StatusBadRequest)
			return
		}	

		if existingCode, ok := store.GetByOriginal(req.URL); ok {
			shortURL := baseURL + "/" + existingCode
			response := map[string]string{"short_url": shortURL}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		code := utils.GenerateShortCode()

		for store.ExistCode(code) {
			code = utils.GenerateShortCode()
		}

		store.Save(code, req.URL)

		shortURL := baseURL + "/" + code
		response := map[string]string{"short_url": shortURL}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}



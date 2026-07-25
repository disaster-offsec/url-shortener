package config

import (
	"os"
)

type Config struct {
	Port string
	BaseURL string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + port
	}

	return &Config{
		Port: port,
		BaseURL: baseURL,
	}
	
}

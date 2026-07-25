package utils

import (
	"crypto/rand"
	"encoding/base64"
	"net/url"
)


// Создаем криптографически стойкий 6-значный URL-безопасный-код
func GenerateShortCode() string {
	b := make([]byte, 4)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:6]
}

// Проверяет на кореектность ссылку
func IsValidURL(raw string) bool {
	parsed, err := url.Parse(raw)
	return err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https")
}

package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var urlMap = make(map[string]string)

func generateShortID(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		originalURL := strings.TrimSpace(string(body))
		if originalURL == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		shortID := generateShortID(8)

		urlMap[shortID] = originalURL

		shortURL := fmt.Sprintf("http://localhost:8080/%s", shortID)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(shortURL))
		if err != nil {
			return
		}

	case http.MethodGet:
		shortID := strings.TrimPrefix(r.URL.Path, "/")

		originalURL, exists := urlMap[shortID]
		if !exists {
			http.Error(w, "Short URL not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

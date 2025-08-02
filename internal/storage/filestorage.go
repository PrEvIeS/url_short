package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

// URLRecord представляет запись URL в файловом хранилище.
type URLRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileStorage struct {
	mu       sync.RWMutex
	data     map[string]string
	filePath string
}

func NewFileStorage(filePath string) (*FileStorage, error) {
	storage := &FileStorage{
		filePath: filePath,
		data:     make(map[string]string),
	}

	if err := storage.load(); err != nil {
		return nil, fmt.Errorf("failed to load storage: %w", err)
	}
	return storage, nil
}

func (s *FileStorage) Set(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	return s.save()
}

func (s *FileStorage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *FileStorage) save() error {
	tmpFile := s.filePath + ".tmp"
	file, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	encoder := json.NewEncoder(file)
	for shortURL, originalURL := range s.data {
		record := URLRecord{
			UUID:        uuid.New().String(),
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		}
		if err := encoder.Encode(record); err != nil {
			return fmt.Errorf("json encoding failed: %w", err)
		}
	}

	if err := os.Rename(tmpFile, s.filePath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}
	return nil
}

func (s *FileStorage) load() error {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to open storage file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	decoder := json.NewDecoder(file)
	for decoder.More() {
		var record URLRecord
		if err := decoder.Decode(&record); err != nil {
			return fmt.Errorf("json decode error: %w", err)
		}
		s.data[record.ShortURL] = record.OriginalURL
	}
	return nil
}

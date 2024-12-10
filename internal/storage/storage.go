package storage

type Storage interface {
	Set(key, value string) error
	Get(key string) (string, bool)
}

type InMemoryStorage struct {
	data map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{data: make(map[string]string)}
}

func (s *InMemoryStorage) Set(key, value string) error {
	s.data[key] = value
	return nil
}

func (s *InMemoryStorage) Get(key string) (string, bool) {
	value, ok := s.data[key]
	return value, ok
}

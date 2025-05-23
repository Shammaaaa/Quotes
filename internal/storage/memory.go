package storage

import (
	"Quotes/internal/models"
	"Quotes/pkg/utils"
	"sync"
	"time"
)

type MemoryStorage struct {
	mu     sync.RWMutex
	quotes []models.Quote
	nextID int
}

func New() *MemoryStorage {
	return &MemoryStorage{
		quotes: make([]models.Quote, 0),
		nextID: 1,
	}
}

func (s *MemoryStorage) Create(q models.Quote) (models.Quote, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if q.Author == "" || q.Text == "" {
		return models.Quote{}, ErrInvalidData
	}

	for _, existing := range s.quotes {
		if existing.Author == q.Author && existing.Text == q.Text {
			return models.Quote{}, ErrAlreadyExists
		}
	}

	newQuote := models.Quote{
		ID:        s.nextID,
		Author:    q.Author,
		Text:      q.Text,
		CreatedAt: time.Now(),
	}

	s.quotes = append(s.quotes, newQuote)
	s.nextID++
	return newQuote, nil
}

func (s *MemoryStorage) GetAll() []models.Quote {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Quote, len(s.quotes))
	copy(result, s.quotes)
	return result
}

func (s *MemoryStorage) GetRandom() (models.Quote, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.quotes) == 0 {
		return models.Quote{}, ErrNotFound
	}

	return s.quotes[utils.RandomInt(0, len(s.quotes))], nil
}

func (s *MemoryStorage) GetByAuthor(author string) []models.Quote {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]models.Quote, 0)
	for _, q := range s.quotes {
		if q.Author == author {
			result = append(result, q)
		}
	}
	return result

}

func (s *MemoryStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, q := range s.quotes {
		if q.ID == id {
			s.quotes = append(s.quotes[:i], s.quotes[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

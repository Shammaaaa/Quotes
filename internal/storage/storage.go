package storage

import "Quotes/internal/models"

type Storage interface {
	Create(quote models.Quote) (models.Quote, error)
	GetAll() []models.Quote
	GetRandom() (models.Quote, error)
	GetByAuthor(author string) []models.Quote
	Delete(id int) error
}

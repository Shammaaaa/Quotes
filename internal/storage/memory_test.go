package storage_test

import (
	"Quotes/internal/models"
	"Quotes/internal/storage"
	"fmt"
	"sync"
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	store := storage.New()

	t.Run("Create - новая цитата", func(t *testing.T) {
		quote := models.Quote{
			Author: "Test Author",
			Text:   "Test Quote Text",
		}

		created, err := store.Create(quote)
		if err != nil {
			t.Fatalf("Ошибка при создании: %v", err)
		}

		if created.ID == 0 {
			t.Error("ID не должен быть 0")
		}
		if created.Author != quote.Author {
			t.Errorf("Ожидался автор %q, получен %q", quote.Author, created.Author)
		}
		if created.CreatedAt.IsZero() {
			t.Error("Дата создания не установлена")
		}
	})

	t.Run("Create - ошибка при пустых данных", func(t *testing.T) {
		_, err := store.Create(models.Quote{Author: "", Text: ""})
		if err != storage.ErrInvalidData {
			t.Errorf("Ожидалась ошибка %v, получена %v", storage.ErrInvalidData, err)
		}
	})

	t.Run("Create - ошибка при дубликате", func(t *testing.T) {
		quote := models.Quote{
			Author: "Duplicate Author",
			Text:   "Duplicate Text",
		}
		store.Create(quote)
		_, err := store.Create(quote)
		if err != storage.ErrAlreadyExists {
			t.Errorf("Ожидалась ошибка %v, получена %v", storage.ErrAlreadyExists, err)
		}
	})

	t.Run("GetAll - получение всех цитат", func(t *testing.T) {
		testStore := storage.New()
		testStore.Create(models.Quote{Author: "Author1", Text: "Text1"})
		testStore.Create(models.Quote{Author: "Author2", Text: "Text2"})

		quotes := testStore.GetAll()
		if len(quotes) != 2 {
			t.Errorf("Ожидалось 2 цитаты, получено %d", len(quotes))
		}
	})

	t.Run("GetRandom - случайная цитата", func(t *testing.T) {
		_, err := store.GetRandom()
		if err != nil {
			t.Errorf("Неожиданная ошибка: %v", err)
		}
	})

	t.Run("GetRandom - ошибка если нет цитат", func(t *testing.T) {
		emptyStore := storage.New()
		_, err := emptyStore.GetRandom()
		if err != storage.ErrNotFound {
			t.Errorf("Ожидалась ошибка %v, получена %v", storage.ErrNotFound, err)
		}
	})

	t.Run("GetByAuthor - фильтрация по автору", func(t *testing.T) {
		store.Create(models.Quote{Author: "Filter Author", Text: "Text1"})
		store.Create(models.Quote{Author: "Filter Author", Text: "Text2"})
		store.Create(models.Quote{Author: "Other Author", Text: "Text3"})

		quotes := store.GetByAuthor("Filter Author")
		if len(quotes) != 2 {
			t.Errorf("Ожидалось 2 цитаты, получено %d", len(quotes))
		}
		for _, q := range quotes {
			if q.Author != "Filter Author" {
				t.Errorf("Найден неверный автор: %q", q.Author)
			}
		}
	})

	t.Run("Delete - удаление цитаты", func(t *testing.T) {
		quote, _ := store.Create(models.Quote{Author: "To Delete", Text: "Text"})
		err := store.Delete(quote.ID)
		if err != nil {
			t.Errorf("Ошибка при удалении: %v", err)
		}

		quotes := store.GetAll()
		for _, q := range quotes {
			if q.ID == quote.ID {
				t.Error("Цитата не была удалена")
			}
		}
	})

	t.Run("Delete - ошибка если не найдена", func(t *testing.T) {
		err := store.Delete(9999)
		if err != storage.ErrNotFound {
			t.Errorf("Ожидалась ошибка %v, получена %v", storage.ErrNotFound, err)
		}
	})

	t.Run("Потокобезопасность", func(t *testing.T) {
		parallelStore := storage.New()
		const numQuotes = 100

		// Используем WaitGroup для ожидания завершения горутин
		var wg sync.WaitGroup
		wg.Add(numQuotes)

		// Запускаем горутины
		for i := 0; i < numQuotes; i++ {
			go func(i int) {
				defer wg.Done()
				_, err := parallelStore.Create(models.Quote{
					Author: "Concurrent Author",
					Text:   fmt.Sprintf("Concurrent Text %d", i),
				})
				if err != nil {
					t.Errorf("Ошибка в горутине %d: %v", i, err)
				}
			}(i)
		}

		// Ждем завершения всех горутин
		wg.Wait()

		// Проверяем результаты
		quotes := parallelStore.GetAll()
		if len(quotes) != numQuotes {
			t.Errorf("Ожидалось %d цитат, получено %d", numQuotes, len(quotes))
		}

		// Проверяем уникальность ID
		ids := make(map[int]bool)
		for _, q := range quotes {
			if ids[q.ID] {
				t.Errorf("Дубликат ID: %d", q.ID)
			}
			ids[q.ID] = true
		}
	})
}

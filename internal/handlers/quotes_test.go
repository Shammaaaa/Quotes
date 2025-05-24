package handlers_test

import (
	"Quotes/internal/handlers"
	"Quotes/internal/models"
	"Quotes/internal/storage"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockStorage struct {
	quotes []models.Quote
}

func (m *MockStorage) Create(q models.Quote) (models.Quote, error) {
	q.ID = len(m.quotes) + 1
	m.quotes = append(m.quotes, q)
	return q, nil
}

func (m *MockStorage) GetAll() []models.Quote {
	return m.quotes
}

func (m *MockStorage) GetRandom() (models.Quote, error) {
	if len(m.quotes) == 0 {
		return models.Quote{}, storage.ErrNotFound
	}
	return m.quotes[0], nil // Всегда возвращаем первую для тестов
}

func (m *MockStorage) GetByAuthor(author string) []models.Quote {
	var result []models.Quote
	for _, q := range m.quotes {
		if q.Author == author {
			result = append(result, q)
		}
	}
	return result
}

func (m *MockStorage) Delete(id int) error {
	for i, q := range m.quotes {
		if q.ID == id {
			m.quotes = append(m.quotes[:i], m.quotes[i+1:]...)
			return nil
		}
	}
	return storage.ErrNotFound
}

func TestQuotesHandler(t *testing.T) {
	mockStorage := &MockStorage{}
	handler := handlers.New(mockStorage)
	router := mux.NewRouter()
	handler.Routes(router)

	t.Run("POST /quotes - создание цитаты", func(t *testing.T) {
		quote := models.CreateRequest{
			Author: "Test Author",
			Text:   "Test Quote",
		}
		body, _ := json.Marshal(quote)
		req := httptest.NewRequest("POST", "/quotes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Ожидался статус %d, получен %d", http.StatusCreated, w.Code)
		}

		var response models.Quote
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatal("Ошибка декодирования ответа:", err)
		}

		if response.Author != quote.Author || response.Text != quote.Text {
			t.Errorf("Ожидалась цитата %+v, получена %+v", quote, response)
		}
	})

	t.Run("GET /quotes - получение всех цитат", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/quotes", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Ожидался статус %d, получен %d", http.StatusOK, w.Code)
		}

		var response []models.Quote
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatal("Ошибка декодирования ответа:", err)
		}

		if len(response) != 1 {
			t.Errorf("Ожидалось 1 цитата, получено %d", len(response))
		}
	})

	t.Run("GET /quotes?author=Test+Author - фильтрация по автору", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/quotes?author=Test+Author", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response []models.Quote
		json.NewDecoder(w.Body).Decode(&response)

		if len(response) != 1 {
			t.Errorf("Ожидалось 1 цитата, получено %d", len(response))
		}
	})

	t.Run("GET /quotes/random - случайная цитата", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/quotes/random", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Ожидался статус %d, получен %d", http.StatusOK, w.Code)
		}

		var response models.Quote
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatal("Ошибка декодирования ответа:", err)
		}
	})

	t.Run("DELETE /quotes/1 - удаление цитаты", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/quotes/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Ожидался статус %d, получен %d", http.StatusNoContent, w.Code)
		}

		if len(mockStorage.quotes) != 0 {
			t.Errorf("Ожидалось 0 цитат после удаления, получено %d", len(mockStorage.quotes))
		}
	})

	t.Run("POST /quotes - ошибка при невалидных данных", func(t *testing.T) {
		testCases := []struct {
			name        string
			body        string
			expectCode  int
			expectError string
		}{
			{
				"Empty author",
				`{"author": "", "text": "text"}`,
				http.StatusUnprocessableEntity,
				"Author and text are required",
			},
			{
				"Empty text",
				`{"author": "author", "text": ""}`,
				http.StatusUnprocessableEntity,
				"Author and text are required",
			},
			{
				"Invalid JSON",
				`{invalid}`,
				http.StatusBadRequest,
				"Invalid request body",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest("POST", "/quotes", bytes.NewReader([]byte(tc.body)))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				if w.Code != tc.expectCode {
					t.Errorf("Ожидался статус %d, получен %d", tc.expectCode, w.Code)
				}

				if tc.expectError != "" {
					var response map[string]interface{}
					if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
						t.Fatal("Ошибка декодирования ответа:", err)
					}
					if response["message"] != tc.expectError {
						t.Errorf("Ожидалась ошибка '%s', получена '%s'",
							tc.expectError, response["message"])
					}
				}
			})
		}
	})

	t.Run("GET /quotes/random - ошибка при отсутствии цитат", func(t *testing.T) {

		mockStorage.quotes = []models.Quote{}

		req := httptest.NewRequest("GET", "/quotes/random", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Ожидался статус %d, получен %d", http.StatusNotFound, w.Code)
		}
	})
}

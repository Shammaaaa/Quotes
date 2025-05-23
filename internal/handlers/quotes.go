package handlers

import (
	"Quotes/internal/models"
	"Quotes/internal/models/api"
	"Quotes/internal/storage"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type QuotesHandler struct {
	store storage.Storage
}

func New(store storage.Storage) *QuotesHandler {
	return &QuotesHandler{store: store}
}

func (h *QuotesHandler) Routes(router *mux.Router) {
	router.HandleFunc("/quotes", h.HandleQuotes).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/quotes/random", h.handleRandomQuote).Methods(http.MethodGet)
	router.HandleFunc("/quotes/{id}", h.handleDeleteQuote).Methods(http.MethodDelete)
}

func (h *QuotesHandler) HandleQuotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listQuotes(w, r)
	case http.MethodPost:
		h.createQuote(w, r)
	default:
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", api.ErrCodeInvalidRequest)

	}
}

func (h *QuotesHandler) listQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	var quotes []models.Quote
	if author != "" {
		quotes = h.store.GetByAuthor(author)
	} else {
		quotes = h.store.GetAll()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

func (h *QuotesHandler) createQuote(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body", api.ErrCodeInvalidRequest)
		return
	}

	if req.Author == "" || req.Text == "" {
		WriteError(w, http.StatusUnprocessableEntity, "Author and text are required", api.ErrCodeInvalidRequest)
		return
	}

	quote, err := h.store.Create(models.Quote{
		Author: req.Author,
		Text:   req.Text,
	})

	if err != nil {
		status := http.StatusInternalServerError
		if err == storage.ErrInvalidData || err == storage.ErrAlreadyExists {
			status = http.StatusUnprocessableEntity

		}
		WriteError(w, status, err.Error(), api.ErrCodeInvalidRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(quote)
}

func (h *QuotesHandler) handleRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := h.store.GetRandom()
	if err != nil {
		if err == storage.ErrNotFound {
			WriteError(w, http.StatusNotFound, "Quote not found", api.ErrCodeNotFound)
		} else {
			WriteError(w, http.StatusInternalServerError, "Internal server error", api.ErrCodeInternal)

		}
		return

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

func (h *QuotesHandler) handleDeleteQuote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid quote ID", api.ErrCodeInvalidRequest)
		return
	}
	if err := h.store.Delete(id); err != nil {
		if err == storage.ErrNotFound {
			WriteError(w, http.StatusNotFound, "Quote not found", api.ErrCodeNotFound)

		} else {
			WriteError(w, http.StatusInternalServerError, "Internal server error", api.ErrCodeInternal)

		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

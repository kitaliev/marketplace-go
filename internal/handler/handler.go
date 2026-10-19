package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/kitaliev/marketplace-go/internal/model"
)

type DataRepository[T any] interface {
	GetByID(context.Context, int64) (T, error)
}

type Handler struct {
	products   DataRepository[model.Product]
	categories DataRepository[model.Category]
}

func New(products DataRepository[model.Product], categories DataRepository[model.Category]) *Handler {
	return &Handler{products: products, categories: categories}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /products/{id}", h.GetProduct)
	mux.HandleFunc("GET /categories/{id}", h.GetCategory)
	return mux
}

func HandleRequestID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return 0, false
	}
	return id, true
}

func HandleMemoryResult(w http.ResponseWriter, name string, err error) {
	if errors.Is(err, model.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": name + " not found"})
		return
	}
	slog.Error("error when returning "+name, "err", err)
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}

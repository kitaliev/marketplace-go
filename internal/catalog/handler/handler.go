package handler

import (
	"context"
	"net/http"

	"github.com/kitaliev/marketplace-go/internal/catalog/model"
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

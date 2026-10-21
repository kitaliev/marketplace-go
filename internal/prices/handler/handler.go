package handler

import (
	"context"
	"net/http"

	"github.com/kitaliev/marketplace-go/internal/prices/model"
)

type DataRepository[T any] interface {
	GetByID(context.Context, int64) (T, error)
}

type Handler struct {
	prices DataRepository[model.Price]
}

func New(prices DataRepository[model.Price]) *Handler {
	return &Handler{prices: prices}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /prices/{id}", h.GetPrices)
	return mux
}

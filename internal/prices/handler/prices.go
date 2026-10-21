package handler

import (
	"net/http"

	"github.com/kitaliev/marketplace-go/internal/platform/httpjson"
)

func (h *Handler) GetPrices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, ok := httpjson.HandleRequestID(w, r)
	if !ok {
		return
	}
	prices, err := h.prices.GetByID(r.Context(), id)
	if err != nil {
		httpjson.HandleMemoryResult(w, "price", err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, prices)
}

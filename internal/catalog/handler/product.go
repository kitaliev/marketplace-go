package handler

import (
	"net/http"

	"github.com/kitaliev/marketplace-go/internal/platform/httpjson"
)

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, ok := httpjson.HandleRequestID(w, r)
	if !ok {
		return
	}
	product, err := h.products.GetByID(r.Context(), id)
	if err != nil {
		httpjson.HandleMemoryResult(w, "product", err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, product)
}

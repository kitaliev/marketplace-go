package handler

import (
	"net/http"
)

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, ok := HandleRequestID(w, r)
	if !ok {
		return
	}
	product, err := h.products.GetByID(r.Context(), id)
	if err != nil {
		HandleMemoryResult(w, "product", err)
		return
	}
	writeJSON(w, http.StatusOK, product)
}

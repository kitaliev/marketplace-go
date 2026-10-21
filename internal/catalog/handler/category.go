package handler

import (
	"net/http"

	"github.com/kitaliev/marketplace-go/internal/platform/httpjson"
)

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, ok := httpjson.HandleRequestID(w, r)
	if !ok {
		return
	}
	category, err := h.categories.GetByID(r.Context(), id)
	if err != nil {
		httpjson.HandleMemoryResult(w, "category", err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, category)
}

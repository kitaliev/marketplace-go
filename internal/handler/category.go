package handler

import (
	"net/http"
)

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, ok := HandleRequestID(w, r)
	if !ok {
		return
	}
	category, err := h.categories.GetByID(r.Context(), id)
	if err != nil {
		HandleMemoryResult(w, "category", err)
		return
	}
	writeJSON(w, http.StatusOK, category)
}

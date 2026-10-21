package httpjson

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/kitaliev/marketplace-go/internal/platform/storage/memory"
)

func HandleRequestID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return 0, false
	}
	return id, true
}

func HandleMemoryResult(w http.ResponseWriter, name string, err error) {
	if errors.Is(err, memory.ErrNotFound) {
		WriteJSON(w, http.StatusNotFound, map[string]string{"error": name + " not found"})
		return
	}
	slog.Error("error when returning "+name, "err", err)
	WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("response failed", "err", err)
	}
}

package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/kitaliev/marketplace-go/internal/platform/storage/memory"

	"github.com/kitaliev/marketplace-go/internal/prices/model"
)

type fakeRepo[T any] struct {
	mu    sync.Mutex
	items map[int64]T
	fail  bool
}

func (r *fakeRepo[T]) GetByID(ctx context.Context, id int64) (T, error) {
	var t T
	if r.fail {
		return t, errors.New("fake error")
	}
	if err := ctx.Err(); err != nil {
		return t, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.items[id]
	if !ok {
		return t, memory.ErrNotFound
	}
	return p, nil
}

var prices = &fakeRepo[model.Price]{
	items: map[int64]model.Price{
		1: {ProductID: 1, Currency: "RUB", Amount: 3000000},
	},
}

func TestGETCategoryAndProducts(t *testing.T) {
	prices.fail = false
	handler := New(prices)
	mux := handler.Routes()
	tests := []struct {
		name       string
		in         string
		wantStatus int
		wantBody   string
	}{
		{"prices: found", "/prices/1", http.StatusOK, `{"product_id":1,"currency":"RUB","amount":3000000}`},
		{"prices: not a number", "/prices/first", http.StatusBadRequest, `{"error":"invalid id"}`},
		{"prices: negative number", "/prices/-1", http.StatusBadRequest, `{"error":"invalid id"}`},
		{"prices: fraction number", "/prices/1.5", http.StatusBadRequest, `{"error":"invalid id"}`},
		{"prices: overflow", "/prices/99999999999999999999", http.StatusBadRequest, `{"error":"invalid id"}`},
		{"prices: zero", "/prices/0", http.StatusBadRequest, `{"error":"invalid id"}`},
		{"prices: not found", "/prices/99", http.StatusNotFound, `{"error":"price not found"}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequestWithContext(context.Background(), http.MethodGet, tc.in, nil)
			mux.ServeHTTP(w, r)
			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("%s - expected Content-Type application/json; got %s", tc.in, w.Header().Get("Content-Type"))
			}
			if w.Code != tc.wantStatus {
				t.Errorf("%s - expected status %d; got %d", tc.in, tc.wantStatus, w.Code)
			}
			if strings.TrimSpace(w.Body.String()) != tc.wantBody {
				t.Errorf("%s - expected body %q; got %q", tc.in, tc.wantBody, w.Body.String())
			}
		})
	}
}

func TestGETCategoryAndProductServerError(t *testing.T) {
	prices.fail = true
	handler := New(prices)
	mux := handler.Routes()

	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/prices/1", nil)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("/prices/1 - extected status 500, got status %d", w.Code)
	}
}

func TestPOSTCategoryAndProductServerError(t *testing.T) {
	prices.fail = false
	handler := New(prices)
	mux := handler.Routes()

	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/prices/1", nil)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("/prices/1 - extected status 405, got status %d", w.Code)
	}
	if w.Header().Get("Allow") != "GET, HEAD" {
		t.Errorf("/prices/1 - extected Header Allow: GET, HEAD, got %s", w.Header().Get("Allow"))
	}
}

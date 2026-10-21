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

	"github.com/kitaliev/marketplace-go/internal/catalog/model"
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

var products = &fakeRepo[model.Product]{
	items: map[int64]model.Product{
		1: {ID: 1, Name: "Test1", Category: "Test Cat", Brand: "TestBrand", SKU: "Test-SKU1"},
		2: {ID: 2, Name: "Test2", Category: "Test Cat", Brand: "TestBrand", SKU: "Test-SKU2"},
		3: {ID: 3, Name: "Test3", Category: "Test Cat", Brand: "TestBrand", SKU: "Test-SKU3"},
	},
}

var categories = &fakeRepo[model.Category]{
	items: map[int64]model.Category{
		1: {ID: 1, Name: "Test1"},
		2: {ID: 2, Name: "Test2"},
		3: {ID: 3, Name: "Test3"},
	},
}

func TestGETCategoryAndProducts(t *testing.T) {
	products.fail = false
	categories.fail = false
	handler := New(products, categories)
	mux := handler.Routes()
	tests := []struct {
		name       string
		in         string
		wantStatus int
		wantBody   string
	}{
		{"categories: found", "/categories/1", http.StatusOK, `{"id":1,"name":"Test1"}`},
		{"categories: not a number", "/categories/first", http.StatusBadRequest, `{"error":"invalid id"}`},
		{"categories: negative number", "/categories/-1", http.StatusBadRequest, `{"error":"invalid id"}`},
		{"categories: fraction number", "/categories/1.5", http.StatusBadRequest, `{"error":"invalid id"}`},
		{"categories: overflow", "/categories/99999999999999999999", http.StatusBadRequest, `{"error":"invalid id"}`},
		{"categories: zero", "/categories/0", http.StatusBadRequest, `{"error":"invalid id"}`},
		{"categories: not found", "/categories/99", http.StatusNotFound, `{"error":"category not found"}`},
		{"products: found", "/products/1", http.StatusOK, `{"id":1,"sku":"Test-SKU1","name":"Test1","category":"Test Cat","brand":"TestBrand"}`},
		{"products: not a number", "/products/first", http.StatusBadRequest, `{"error":"invalid id"}`},
		{"products: negative number", "/products/-1", http.StatusBadRequest, `{"error":"invalid id"}`},
		{"products: fraction number", "/products/1.5", http.StatusBadRequest, `{"error":"invalid id"}`},
		{"products: overflow", "/products/99999999999999999999", http.StatusBadRequest, `{"error":"invalid id"}`},
		{"products: zero", "/products/0", http.StatusBadRequest, `{"error":"invalid id"}`},
		{"products: not found", "/products/99", http.StatusNotFound, `{"error":"product not found"}`},
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
	products.fail = true
	categories.fail = true
	handler := New(products, categories)
	mux := handler.Routes()

	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/categories/1", nil)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("/categories/1 - extected status 500, got status %d", w.Code)
	}

	w = httptest.NewRecorder()
	r = httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/products/1", nil)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("/products/1 - extected status 500, got status %d", w.Code)
	}
}

func TestPOSTCategoryAndProductServerError(t *testing.T) {
	products.fail = false
	categories.fail = false
	handler := New(products, categories)
	mux := handler.Routes()

	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/categories/1", nil)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("/categories/1 - extected status 405, got status %d", w.Code)
	}
	if w.Header().Get("Allow") != "GET, HEAD" {
		t.Errorf("/categories/1 - extected Header Allow: GET, HEAD, got %s", w.Header().Get("Allow"))
	}

	w = httptest.NewRecorder()
	r = httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/products/1", nil)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("/products/1 - extected status 405, got status %d", w.Code)
	}
	if w.Header().Get("Allow") != "GET, HEAD" {
		t.Errorf("/products/1 - extected Header Allow: GET, HEAD, got %s", w.Header().Get("Allow"))
	}
}

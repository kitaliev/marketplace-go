package memory

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/kitaliev/marketplace-go/internal/seed"

	"go.uber.org/goleak"

	"github.com/kitaliev/marketplace-go/internal/model"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestNew(t *testing.T) {
	list := seed.ParsedList{
		Products: []model.Product{
			{ID: 1, Name: "Test", Category: "Test", SKU: "Test", Brand: "Test"},
			{ID: 2, Name: "Test", Category: "Test", SKU: "Test", Brand: "Test"},
		},
		Categories: []model.Category{
			{ID: 1, Name: "Test"},
			{ID: 2, Name: "Test"},
		},
	}
	t.Run("products", func(t *testing.T) {
		got := New(list.Products)
		want := map[int64]model.Product{
			1: list.Products[0],
			2: list.Products[1],
		}
		if !reflect.DeepEqual(got.items, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("categories", func(t *testing.T) {
		got := New(list.Categories)
		want := map[int64]model.Category{
			1: list.Categories[0],
			2: list.Categories[1],
		}
		if !reflect.DeepEqual(got.items, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestGetByID(t *testing.T) {
	list := seed.ParsedList{
		Products: []model.Product{
			{ID: 1, Name: "Test", Category: "Test", SKU: "Test", Brand: "Test"},
			{ID: 2, Name: "Test", Category: "Test", SKU: "Test", Brand: "Test"},
		},
		Categories: []model.Category{
			{ID: 1, Name: "Test"},
			{ID: 2, Name: "Test"},
		},
	}
	tests := []struct {
		name    string
		in      int64
		want    model.Product
		wantErr error
	}{
		{"normal case", 1, list.Products[0], nil},
		{"wrong id", -1, model.Product{}, model.ErrNotFound},
		{"no product", 999, model.Product{}, model.ErrNotFound},
	}
	store := New(list.Products)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := store.GetByID(context.Background(), tc.in)
			if got != tc.want {
				t.Errorf("GetByID(%d) = %v; want %v", tc.in, got, tc.want)
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("GetByID(%d) got error %v; expected error %v", tc.in, err, tc.wantErr)
			}
		})
	}
}

func TestGetByIDNotChangingOriginal(t *testing.T) {
	list := seed.ParsedList{
		Products: []model.Product{
			{ID: 1, Name: "Test", Category: "Test", SKU: "Test", Brand: "Test"},
			{ID: 2, Name: "Test", Category: "Test", SKU: "Test", Brand: "Test"},
		},
		Categories: []model.Category{
			{ID: 1, Name: "Test"},
			{ID: 2, Name: "Test"},
		},
	}
	store := New(list.Products)
	got, err := store.GetByID(context.Background(), 1)
	original := got
	if err != nil {
		t.Fatalf(`GetByID(1) from seed_test failed with error %v`, err)
	}
	got.ID = 111
	got.SKU = "TEST-SKU"
	got.Name = "TEST"
	got.Category = "TEST CAT"
	got.Brand = "TEST Brand"

	got2, err := store.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf(`GetByID(1) (got2) from seed_test failed with error %v`, err)
	}
	if got2 != original {
		t.Errorf("Changing data that was received from GetByID(1) changed data in list which should not happen")
	}
}

func TestSave(t *testing.T) {
	list := seed.ParsedList{
		Products: []model.Product{
			{ID: 1, Name: "Test", Category: "Test", SKU: "Test", Brand: "Test"},
			{ID: 2, Name: "Test", Category: "Test", SKU: "Test", Brand: "Test"},
		},
		Categories: []model.Category{
			{ID: 1, Name: "Test"},
			{ID: 2, Name: "Test"},
		},
	}
	newProduct := model.Product{
		ID:       11,
		SKU:      "NEW-SKU",
		Name:     "NEW PRODUCT",
		Category: "NEW CATEGORY",
		Brand:    "NEW BRAND",
	}
	existingProduct := model.Product{
		ID:       1,
		SKU:      "RUDNC1002-STHBCH",
		Name:     "Rudis Journey Knit Youth",
		Category: "TEST",
		Brand:    "Rudis",
	}

	store := New(list.Products)
	err := store.Save(context.Background(), newProduct)
	if err != nil {
		t.Errorf("Save returned error %v, expected no error", err)
	}
	got, err := store.GetByID(context.Background(), 11)
	if err != nil {
		t.Errorf("GetByID(11) after save should return it, error received %v", err)
	}
	if got != newProduct {
		t.Errorf("GetByID(11) after save should returned something else, expected %v, got %v", newProduct, got)
	}

	err = store.Save(context.Background(), existingProduct)
	if err != nil {
		t.Errorf("Save returned error %v, expected no error", err)
	}

	got2, err := store.GetByID(context.Background(), 1)
	if err != nil {
		t.Errorf("GetByID(1) after save returned error %v; expected no error", err)
	}
	if got2 != existingProduct {
		t.Errorf("GetByID(1) after save returned old product; ether changes was not applied or method returns old data")
	}
}

func TestLeakNLoss(t *testing.T) {
	list := seed.ParsedList{
		Products: []model.Product{
			{ID: 1, Name: "Test", Category: "Test", SKU: "Test", Brand: "Test"},
			{ID: 2, Name: "Test", Category: "Test", SKU: "Test", Brand: "Test"},
		},
		Categories: []model.Category{
			{ID: 1, Name: "Test"},
			{ID: 2, Name: "Test"},
		},
	}
	store := New(list.Products)
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Go(func() {
			err := store.Save(context.Background(), model.Product{
				ID:       int64(i + 100),
				SKU:      "SKU-TEST",
				Name:     "Test",
				Category: "Test Category",
				Brand:    "Test Brand",
			})
			if err != nil {
				t.Errorf("Save returned error %v, expected no error", err)
			}
			v, err := store.GetByID(context.Background(), int64(i+100))
			if err != nil {
				t.Errorf("GetByID(%d) returned error %v, no error expected", i+100, err)
			} else if v.ID != int64(i+100) {
				t.Errorf("GetByID(%d) returned product with id %d, expected %d", i+100, v.ID, i+100)
			}
		})
	}
	wg.Wait()
}

package jsonfile

import (
	"errors"
	"sync"
	"testing"

	"go.uber.org/goleak"

	"github.com/kitaliev/marketplace-go/internal/model"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestLoadAndValidateProducts(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		wantData bool
		wantErr  bool
	}{
		{"normal file", "testdata/seed_test.json", true, false},
		{"no file", "testdata/seed_test_123.json", false, true},
		{"duplicate id", "testdata/seed_test_duplicate_id.json", false, true},
		{"empty file", "testdata/seed_test_empty.json", false, true},
		{"empty category", "testdata/seed_test_empty_category.json", false, true},
		{"empty name", "testdata/seed_test_empty_name.json", false, true},
		{"empty sku", "testdata/seed_test_empty_sku.json", false, true},
		{"negative id", "testdata/seed_test_negative_id.json", false, true},
		{"zero id", "testdata/seed_test_zero_id.json", false, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := Load(tc.in)
			gotErr := err != nil
			gotData := data != nil
			if gotErr != tc.wantErr {
				t.Errorf("Load() - expected error? - %t, got error? - %t: %v", tc.wantErr, gotErr, gotErr)
			}
			if gotData != tc.wantData {
				t.Errorf("Load() expected to be data? - %t, got data? - %t", tc.wantData, gotData)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	product := model.Product{
		ID:       1,
		SKU:      "RUDNC1002-STHBCH",
		Name:     "Rudis Journey Knit Youth",
		Category: "Кроссовки",
		Brand:    "Rudis",
	}
	tests := []struct {
		name    string
		in      int64
		want    model.Product
		wantErr error
	}{
		{"normal case", 1, product, nil},
		{"wrong id", -1, model.Product{}, ErrIncorrectID},
		{"no product", 999, model.Product{}, ErrProductNotFound},
	}
	f, err := Load("testdata/seed_test.json")
	if err != nil {
		t.Fatalf(`Load("testdata/seed_test.json") failed with error %v`, err)
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := f.GetByID(tc.in)
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
	f, err := Load("testdata/seed_test.json")
	if err != nil {
		t.Fatalf(`Load("testdata/seed_test.json") failed with error %v`, err)
	}
	got, err := f.GetByID(1)
	original := got
	if err != nil {
		t.Fatalf(`GetByID(1) from seed_test failed with error %v`, err)
	}
	got.ID = 111
	got.SKU = "TEST-SKU"
	got.Name = "TEST"
	got.Category = "TEST CAT"
	got.Brand = "TEST Brand"

	got2, err := f.GetByID(1)
	if err != nil {
		t.Fatalf(`GetByID(1) (got2) from seed_test failed with error %v`, err)
	}
	if got2 != original {
		t.Errorf("Changing data that was received from GetByID(1) changed data in list which should not happen")
	}
}

func TestSave(t *testing.T) {
	f, err := Load("testdata/seed_test.json")
	if err != nil {
		t.Fatalf(`Load("testdata/seed_test.json") failed with error %v`, err)
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

	f.Save(newProduct)
	got, err := f.GetByID(11)
	if err != nil {
		t.Errorf("GetByID(11) after save should return it, error received %v", err)
	}
	if got != newProduct {
		t.Errorf("GetByID(11) after save should returned something else, expected %v, got %v", newProduct, got)
	}

	f.Save(existingProduct)

	got2, err := f.GetByID(1)
	if err != nil {
		t.Errorf("GetByID(1) after save returned error %v; expected no error", err)
	}
	if got2 != existingProduct {
		t.Errorf("GetByID(1) after save returned old product; ether changes was not applied or method returns old data")
	}
}

func TestLeakNLoss(t *testing.T) {
	f, err := Load("testdata/seed_test.json")
	if err != nil {
		t.Fatalf(`Load("testdata/seed_test.json") failed with error %v`, err)
	}

	var wg sync.WaitGroup
	for i := range 50 {
		wg.Go(func() {
			f.Save(model.Product{
				ID:       int64(i + 100),
				SKU:      "SKU-TEST",
				Name:     "Test",
				Category: "Test Category",
				Brand:    "Test Brand",
			})
			v, err := f.GetByID(int64(i + 100))
			if err != nil {
				t.Errorf("GetByID(%d) returned error %v, no error expected", i+100, err)
			} else if v.ID != int64(i+100) {
				t.Errorf("GetByID(%d) returned product with id %d, expected %d", i+100, v.ID, i+100)
			}
		})
	}
	wg.Wait()
}

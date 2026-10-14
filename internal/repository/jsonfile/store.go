package jsonfile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/kitaliev/marketplace-go/internal/model"
)

var ErrProductNotFound = errors.New("product not found")

type File struct {
	mu       sync.Mutex
	products map[int64]model.Product
}

func Load(fileName string) (*File, error) {
	data, readErr := os.ReadFile(fileName)
	if readErr != nil {
		return nil, fmt.Errorf("unable to read file: %w", readErr)
	}

	var placeholder struct {
		Products []model.Product
	}
	parseErr := json.Unmarshal(data, &placeholder)
	if parseErr != nil {
		return nil, fmt.Errorf("unable to parse file: %w", parseErr)
	}
	validateErr := validateProducts(placeholder.Products)
	if validateErr != nil {
		return nil, fmt.Errorf("error happened when validating products: %w", validateErr)
	}

	var f File
	f.products = make(map[int64]model.Product, len(placeholder.Products))
	for _, p := range placeholder.Products {
		f.products[p.ID] = p
	}
	return &f, nil
}

func validateProducts(products []model.Product) error {
	seedIDs := make(map[int64]bool, len(products))
	for _, p := range products {
		if seedIDs[p.ID] {
			return fmt.Errorf("multiple products have same ids: %v", p)
		}
		if p.ID <= 0 {
			return fmt.Errorf("product id is less than or equals to zero: %v", p)
		}
		if p.SKU == "" {
			return fmt.Errorf("product have empty SKU field: %v", p)
		}
		if p.Name == "" {
			return fmt.Errorf("product have empty Name field: %v", p)
		}
		if p.Category == "" {
			return fmt.Errorf("product have empty Category field: %v", p)
		}
		seedIDs[p.ID] = true
	}
	return nil
}

func (f *File) GetByID(id int64) (model.Product, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	p, ok := f.products[id]
	if !ok {
		return model.Product{}, ErrProductNotFound
	}
	return p, nil
}

func (f *File) Save(product model.Product) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.products[product.ID] = product
}

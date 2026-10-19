package model

import (
	"fmt"
)

type Product struct {
	ID       int64  `json:"id"`
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Brand    string `json:"brand"`
}

func (p Product) GetID() int64 {
	return p.ID
}

func ValidateProducts(items []Product) error {
	seedIDs := make(map[int64]bool, len(items))
	for _, p := range items {
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

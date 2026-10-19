package model

import (
	"fmt"
)

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (c Category) GetID() int64 {
	return c.ID
}

func ValidateCategories(items []Category) error {
	seedIDs := make(map[int64]bool, len(items))
	for _, c := range items {
		if seedIDs[c.ID] {
			return fmt.Errorf("multiple categories have same ids: %v", c)
		}
		if c.ID <= 0 {
			return fmt.Errorf("categories id is less than or equals to zero: %v", c)
		}
		if c.Name == "" {
			return fmt.Errorf("categories have empty Name field: %v", c)
		}
		seedIDs[c.ID] = true
	}
	return nil
}

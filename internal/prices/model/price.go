package model

import (
	"fmt"
	"strings"
)

type Price struct {
	ProductID int64  `json:"product_id"`
	Currency  string `json:"currency"`
	Amount    int64  `json:"amount"`
}

func (c Price) GetID() int64 {
	return c.ProductID
}

func ValidatePrices(items []Price) error {
	seedIDs := make(map[int64]bool, len(items))
	for _, c := range items {
		if seedIDs[c.ProductID] {
			return fmt.Errorf("multiple prices have same product id: %v", c)
		}
		if c.ProductID <= 0 {
			return fmt.Errorf("product id field should be greater than 0: %v", c)
		}
		if c.Currency == "" {
			return fmt.Errorf("currency field is empty: %v", c)
		}
		if len(c.Currency) != 3 {
			return fmt.Errorf("currency field should contain 3 letters only e.g. RUB: %v", c)
		}
		if strings.ToUpper(c.Currency) != c.Currency {
			return fmt.Errorf("currency field should contain 3 letters in uppercase e.g. RUB: %v", c)
		}
		if c.Amount <= 0 {
			return fmt.Errorf("amount field should be greater than 0: %v", c)
		}
		seedIDs[c.ProductID] = true
	}
	return nil
}

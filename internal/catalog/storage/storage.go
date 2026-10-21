package storage

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/kitaliev/marketplace-go/internal/catalog/model"
)

type ProductDTO struct {
	ID       int64  `json:"id"`
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Brand    string `json:"brand"`
}

type CategoryDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type DTOList struct {
	Products   []ProductDTO  `json:"products"`
	Categories []CategoryDTO `json:"categories"`
}

type ParsedList struct {
	Products   []model.Product
	Categories []model.Category
}

func Read(fileName string) (*ParsedList, error) {
	data, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := data.Close()
		if err != nil {
			slog.Error("error closing file", "err", err)
		}
	}()

	parsed := DTOList{}
	dec := json.NewDecoder(data)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&parsed); err != nil {
		return nil, err
	}

	products := make([]model.Product, 0, len(parsed.Products))
	for _, p := range parsed.Products {
		product := model.Product{
			ID:       p.ID,
			SKU:      p.SKU,
			Name:     p.Name,
			Category: p.Category,
			Brand:    p.Brand,
		}
		products = append(products, product)
	}
	if err := model.ValidateProducts(products); err != nil {
		return nil, err
	}

	categories := make([]model.Category, 0, len(parsed.Categories))
	for _, p := range parsed.Categories {
		category := model.Category{
			ID:   p.ID,
			Name: p.Name,
		}
		categories = append(categories, category)
	}
	if err := model.ValidateCategories(categories); err != nil {
		return nil, err
	}
	result := ParsedList{
		Products:   products,
		Categories: categories,
	}
	return &result, nil
}

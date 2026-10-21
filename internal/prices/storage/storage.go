package storage

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/kitaliev/marketplace-go/internal/prices/model"
)

type PriceDTO struct {
	ProductID int64  `json:"product_id"`
	Currency  string `json:"currency"`
	Amount    int64  `json:"amount"`
}

type DTOList struct {
	Prices []PriceDTO `json:"prices"`
}

type ParsedList struct {
	Prices []model.Price
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

	prices := make([]model.Price, 0, len(parsed.Prices))
	for _, p := range parsed.Prices {
		price := model.Price{
			ProductID: p.ProductID,
			Currency:  p.Currency,
			Amount:    p.Amount,
		}
		prices = append(prices, price)
	}
	if err := model.ValidatePrices(prices); err != nil {
		return nil, err
	}
	result := ParsedList{
		Prices: prices,
	}
	return &result, nil
}

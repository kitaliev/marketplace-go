package model

import (
	"testing"
)

func TestValidateCategories(t *testing.T) {
	tests := []struct {
		name    string
		in      []Price
		wantErr bool
	}{
		{"normal case", []Price{{1, "RUB", 300000}, {2, "RUB", 300000}}, false},
		{"duplicate id", []Price{{1, "RUB", 300000}, {1, "RUB", 300000}}, true},
		{"empty currency", []Price{{1, "", 300000}}, true},
		{"wrong currency", []Price{{1, "$$$$$", 300000}}, true},
		{"not uppercase currency", []Price{{1, "rub", 300000}}, true},
		{"negative id", []Price{{-1, "RUB", 300000}}, true},
		{"zero id", []Price{{0, "RUB", 300000}}, true},
		{"negative amount", []Price{{1, "rub", -1}}, true},
		{"zero amount", []Price{{1, "rub", 0}}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePrices(tc.in)
			gotErr := err != nil
			if gotErr != tc.wantErr {
				t.Errorf("ValidatePrices(%v) - expected error? - %t, got error? - %t: %v", tc.in, tc.wantErr, gotErr, err)
			}
		})
	}
}

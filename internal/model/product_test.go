package model

import (
	"testing"
)

func TestValidateProducts(t *testing.T) {
	tests := []struct {
		name    string
		in      []Product
		wantErr bool
	}{
		{
			"normal case",
			[]Product{
				{1, "Test", "Test", "Test", "Test"},
				{2, "Test", "Test", "Test", "Test"},
			},
			false,
		},
		{
			"duplicate id",
			[]Product{
				{1, "Test", "Test", "Test", "Test"},
				{1, "Test", "Test", "Test", "Test"},
			},
			true,
		},
		{
			"empty category",
			[]Product{
				{1, "Test", "Test", "", "Test"},
			},
			true,
		},
		{
			"empty name",
			[]Product{
				{1, "Test", "", "Test", "Test"},
			},
			true,
		},
		{
			"empty sku",
			[]Product{
				{1, "", "Test", "Test", "Test"},
			},
			true,
		},
		{
			"wrong id",
			[]Product{
				{-1, "Test", "Test", "Test", "Test"},
			},
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateProducts(tc.in)
			gotErr := err != nil
			if gotErr != tc.wantErr {
				t.Errorf("ValidateCategories(%v) - expected error? - %t, got error? - %t: %v", tc.in, tc.wantErr, gotErr, err)
			}
		})
	}
}

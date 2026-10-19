package model

import (
	"testing"
)

func TestValidateCategories(t *testing.T) {
	tests := []struct {
		name    string
		in      []Category
		wantErr bool
	}{
		{"normal case", []Category{{1, "Test"}, {2, "OK"}}, false},
		{"duplicate id", []Category{{1, "Test"}, {1, "OK"}}, true},
		{"empty name", []Category{{1, ""}}, true},
		{"negative id", []Category{{-1, "Test"}}, true},
		{"zero id", []Category{{0, "Test"}}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCategories(tc.in)
			gotErr := err != nil
			if gotErr != tc.wantErr {
				t.Errorf("ValidateCategories(%v) - expected error? - %t, got error? - %t: %v", tc.in, tc.wantErr, gotErr, err)
			}
		})
	}
}

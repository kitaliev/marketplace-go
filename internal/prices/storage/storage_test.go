package storage

import (
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestNewAndValidate(t *testing.T) {
	tests := []struct {
		name       string
		in         string
		wantLength int
		wantErr    bool
	}{
		{"normal file", "testdata/seed_test.json", 9, false},
		{"no file", "testdata/seed_test_123.json", 0, true},
		{"empty file", "testdata/seed_test_empty.json", 0, true},
		{"wrong file", "testdata/seed_test_wrong_file.json", 0, true},
		{"wrong data", "testdata/seed_test_wrong_data.json", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := Read(tc.in)
			gotErr := err != nil
			if gotErr != tc.wantErr {
				t.Errorf("Read(%s) - expected error? - %t, got error? - %t: %v", tc.in, tc.wantErr, gotErr, err)
			}
			if err == nil && len(data.Prices) != tc.wantLength {
				t.Errorf("Read(%s) - expected to return %d elements; got %d", tc.in, tc.wantLength, len(data.Prices))
			}
		})
	}
}

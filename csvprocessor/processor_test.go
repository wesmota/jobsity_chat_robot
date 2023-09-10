package csvprocessor

import (
	"os"
	"testing"
)

func TestReadCSVStockFile(t *testing.T) {
	testCases := []struct {
		fileName string
		expected string
	}{
		{
			"correct.csv",
			"AAPL.US quote is $178.18 per share",
		},
		{
			"not_available.csv",
			"AAPL.US quote is not available",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			file, err := os.Open("seeds/" + tc.fileName)
			if err != nil {
				t.Errorf("Error opening file: %v", err)
				return
			}
			defer file.Close()

			result := readCSVStockFile(file)
			if result != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, result)
			}
		})
	}
}

package compare

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestCompareFiles(t *testing.T) {
	// arrange
	expectedResult := PasswordItem{Name: "example.com", Username: "adrian", SHA1: "00000099A4D3034E14DF60EF50799F695C27C0EC"}

	// act
	actualResults, err := CompareFiles("../../test/hibp.txt", "../../test/wpm.json")

	// assert
	if err != nil {
		t.Error(err)
	}

	if len(actualResults) != 1 {
		t.Errorf("Expected 1 match, got %d", len(actualResults))
	}
	if expectedResult != actualResults[0] {
		t.Errorf("Expected result '%v' did not match actual '%v'", expectedResult, actualResults)
	}
}

func BenchmarkCompareFiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := CompareFiles("../../test/hibp.txt", "../../test/wpm.json")
		if err != nil {
			b.Error(err)
		}
	}
}
func TestReadPasswordItems(t *testing.T) {
	// arrange
	jsonData := `
	{
		"passwords": [
			{
				"name": "example.com",
				"username": "adrian",
				"sha1": "64A6DA114D17AE8F167F6BE2C4AEBC9E99F7466C"
			},
			{
				"name": "example.org",
				"username": "thomas",
				"sha1": "B1B3773A05C0ED0176787A4F1574FF0075F7521E"
			}
		]
	}`
	expectedItems := []PasswordItem{
		{
			Name:     "example.com",
			Username: "adrian",
			SHA1:     "64A6DA114D17AE8F167F6BE2C4AEBC9E99F7466C",
		},
		{
			Name:     "example.org",
			Username: "thomas",
			SHA1:     "B1B3773A05C0ED0176787A4F1574FF0075F7521E",
		},
	}
	reader := strings.NewReader(jsonData)

	// act
	passwordItems, err := readPasswordItems(reader)

	// assert
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(passwordItems, expectedItems) {
		t.Errorf("Expected password items to be %v, but got %v", expectedItems, passwordItems)
	}
}
func TestFindHash(t *testing.T) {
	// arrange
	data := []byte("HASH1:1\nHASH2:2\nHASH3:3\nHASH4:4\nHASH5:5\n")
	start := int64(0)
	end := int64(len(data))
	file, err := os.CreateTemp("", "test_file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	_, err = file.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	// arrange
	tests := []struct {
		name          string
		start         int64
		end           int64
		hashToFind    string
		expectedHash  string
		expectedCount int64
		expectedError string
	}{
		{
			name:          "Valid hash found",
			start:         start,
			end:           end,
			hashToFind:    "HASH3",
			expectedHash:  "HASH3",
			expectedCount: 3,
			expectedError: "",
		},
		{
			name:          "Hash not found",
			start:         start,
			end:           end,
			hashToFind:    "INVALID_HASH",
			expectedHash:  "",
			expectedCount: 0,
			expectedError: "",
		},
		{
			name:          "Invalid range",
			start:         end,   // topsy turvy
			end:           start, // topsy turvy
			hashToFind:    "INVALID_RANGE",
			expectedHash:  "",
			expectedCount: 0,
			expectedError: "invalid range",
		},
		{
			name:          "Invalid range",
			start:         -10, // woops
			end:           10,
			hashToFind:    "INVALID_RANGE",
			expectedHash:  "",
			expectedCount: 0,
			expectedError: "invalid range",
		},
		{
			name:          "Invalid range",
			start:         10,
			end:           -10, // woops
			hashToFind:    "INVALID_RANGE",
			expectedHash:  "",
			expectedCount: 0,
			expectedError: "invalid range",
		},
		{
			name:          "Hash at start of range",
			start:         start,
			end:           end,
			hashToFind:    "HASH1",
			expectedHash:  "HASH1",
			expectedCount: 1,
			expectedError: "",
		},
		{
			name:          "Has at end of range",
			start:         start,
			end:           end,
			hashToFind:    "HASH5",
			expectedHash:  "HASH5",
			expectedCount: 5,
			expectedError: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// act
			gotHash, gotCount, err := findHash(test.start, test.end, file, test.hashToFind)

			// assert
			if err != nil && err.Error() != test.expectedError {
				t.Errorf("Unexpected error: %v, expected: %v", err, test.expectedError)
			}
			if gotHash != test.expectedHash {
				t.Errorf("findHash(%d, %d, file, %q) = %q, want %q", test.start, test.end, test.hashToFind, gotHash, test.expectedHash)
			}
			if gotCount != test.expectedCount {
				t.Errorf("findHash(%d, %d, file, %q) count = %d, want %d", test.start, test.end, test.hashToFind, gotCount, test.expectedCount)
			}
		})
	}
}

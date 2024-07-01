package compare

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

// func generateHIBPTestData(size int) []string {
// 	testData := make([]string, size)

// 	for i := 0; i < size; i++ {
// 		countPrefix := rand.Intn(10000) + 1
// 		rawHash := sha1.New()
// 		rawHash.Write([]byte(fmt.Sprint(i)))
// 		hash := hex.EncodeToString(rawHash.Sum(nil)) + ":" + fmt.Sprint(countPrefix)
// 		testData[i] = hash
// 	}
// 	slices.Sort(testData)

// 	return testData
// }

// func generateWPMTestData(size int) []PasswordItem {
// 	testData := make([]PasswordItem, size)
// 	for i := 0; i < size; i++ {
// 		rawHash := sha1.New()
// 		rawHash.Write([]byte(fmt.Sprint(i * i)))
// 		testData[i] = PasswordItem{"example.com" + fmt.Sprint(i), "adrian " + fmt.Sprint(i), hex.EncodeToString(rawHash.Sum(nil))}
// 	}

// 	return testData
// }

// TODO refactor
// func TestCompare(t *testing.T) {
// 	// arrange
// 	HIBP := []string{
// 		"01B307ACBA4F54F55AAFC33BB06BBBF6CA803E9A:100",
// 		"5BAA61E4C9B93F3F0682250B6CF8331B7EE68FD8:20",
// 		"64A6DA114D17AE8F167F6BE2C4AEBC9E99F7466C:1",
// 		"B1B3773A05C0ED0176787A4F1574FF0075F7521E:9001",
// 	}
// 	WPM := []PasswordItem{
// 		{"example.com", "adrian", "64A6DA114D17AE8F167F6BE2C4AEBC9E99F7466C"},
// 	}

// 	// act
// 	matches := Compare(HIBP, WPM)

// 	// assert
// 	if len(matches) != 1 {
// 		t.Errorf("Expected 1 match, got %d", len(matches))
// 	}
// 	match := matches[0]
// 	if match.Name != "example.com" {
// 		t.Errorf("Expected name to be example.com, got %s", match.Name)
// 	}
// 	if match.Username != "adrian" {
// 		t.Errorf("Expected username to be adrian, got %s", match.Username)
// 	}
// 	if match.SHA1 != "64A6DA114D17AE8F167F6BE2C4AEBC9E99F7466C" {
// 		t.Errorf("Expected sha1 to be 64A6DA114D17AE8F167F6BE2C4AEBC9E99F7466C, got %s", match.SHA1)
// 	}
// }

// TODO refactor
// func BenchmarkCompare(b *testing.B) {
// 	for _, testCaseSize := range []struct {
// 		HIBPSize int
// 		WPMSize  int
// 	}{
// 		{HIBPSize: 10, WPMSize: 1000},
// 		{HIBPSize: 100, WPMSize: 1000},
// 		{HIBPSize: 1_000, WPMSize: 1000},
// 		{HIBPSize: 10_000, WPMSize: 1000},
// 		{HIBPSize: 100_000, WPMSize: 1000},
// 		{HIBPSize: 1_000_000, WPMSize: 1000},
// 	} {
// 		HIBP := generateHIBPTestData(testCaseSize.HIBPSize)
// 		WPM := generateWPMTestData(testCaseSize.WPMSize)

// 		if len(HIBP) != testCaseSize.HIBPSize {
// 			panic(fmt.Sprintf("Problem with test data: Requested HIBP length doesn't match. Expected: %v actual: %v", testCaseSize.HIBPSize, len(HIBP)))
// 		}
// 		if len(WPM) != testCaseSize.WPMSize {
// 			panic("Problem with test data: Requested WPM length doesn't match")
// 		}

// 		for i := 0; i < b.N; i++ {
// 			Compare(HIBP, WPM)
// 		}
// 	}
// }

func TestCompareFiles(t *testing.T) {
	// arrange
	expectedResult := PasswordItem{Name: "example.com", Username: "adrian", SHA1: "00000099A4D3034E14DF60EF50799F695C27C0EC"}

	// act
	actualResults, err := CompareFiles("examples/hibp.txt", "examples/wpm.json")

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
		_, err := CompareFiles("examples/hibp.txt", "examples/wpm.json")
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
	hashToFind := "HASH3"
	start := int64(0)
	end := int64(len(data))
	expectedCount := int64(3)
	file, err := os.CreateTemp("", "test_file")
	if err != nil {
		t.Fatal(err)
	}
	// defer os.Remove(file.Name())

	_, err = file.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	// act
	gotHash, gotCount, err := findHash(start, end, file, hashToFind)

	// assert
	// TODO - add a test table?
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if gotHash != hashToFind {
		t.Errorf("findHash(%d, %d, file, %q) = %q, want %q", start, end, hashToFind, gotHash, expectedCount)
	}
	if gotCount != expectedCount {
		t.Errorf("findHash(%d, %d, file, %q) count = %d, want %d", start, end, hashToFind, gotCount, expectedCount)
	}

}

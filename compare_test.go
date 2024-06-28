package compare

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
	"slices"
	"testing"
)

func generateHIBPTestData(size int) []string {
	testData := make([]string, size)

	for i := 0; i < size; i++ {
		countPrefix := rand.Intn(10000) + 1
		rawHash := sha1.New()
		rawHash.Write([]byte(fmt.Sprint(i)))
		hash := hex.EncodeToString(rawHash.Sum(nil)) + ":" + fmt.Sprint(countPrefix)
		testData[i] = hash
	}
	slices.Sort(testData)

	return testData
}

func generateWPMTestData(size int) []PasswordItem {
	testData := make([]PasswordItem, size)
	for i := 0; i < size; i++ {
		rawHash := sha1.New()
		rawHash.Write([]byte(fmt.Sprint(i * i)))
		testData[i] = PasswordItem{"example.com" + fmt.Sprint(i), "adrian " + fmt.Sprint(i), hex.EncodeToString(rawHash.Sum(nil))}
	}

	return testData
}

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
	found, err := CompareFiles("examples/hibp.txt", "examples/wpm.json")

	if err != nil {
		t.Error(err)
	}

	fmt.Println(found[0]) // TODO log.
}

func BenchmarkCompareFiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := CompareFiles("examples/hibp.txt", "examples/wpm.json")
		if err != nil {
			b.Error(err)
		}
	}
}

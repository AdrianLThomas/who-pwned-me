package compare

import (
	"testing"
)

func getHIBPTestData() []string {
	return []string{
		"01B307ACBA4F54F55AAFC33BB06BBBF6CA803E9A:100",
		"5BAA61E4C9B93F3F0682250B6CF8331B7EE68FD8:20",
		"64A6DA114D17AE8F167F6BE2C4AEBC9E99F7466C:1",
		"B1B3773A05C0ED0176787A4F1574FF0075F7521E:9001",
	}
}

func getWPMTestData() []PasswordItem {
	return []PasswordItem{
		{"example.com", "adrian", "64A6DA114D17AE8F167F6BE2C4AEBC9E99F7466C"},
	}
}

func TestCompare(t *testing.T) {
	// arrange
	HIBP := getHIBPTestData()
	WPM := getWPMTestData()

	// act
	matches := Compare(HIBP, WPM)

	// assert
	if len(matches) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matches))
	}
	match := matches[0]
	if match.name != "example.com" {
		t.Errorf("Expected name to be example.com, got %s", match.name)
	}
	if match.username != "adrian" {
		t.Errorf("Expected username to be adrian, got %s", match.username)
	}
	if match.password != "64A6DA114D17AE8F167F6BE2C4AEBC9E99F7466C" {
		t.Errorf("Expected sha1 to be 64A6DA114D17AE8F167F6BE2C4AEBC9E99F7466C, got %s", match.password)
	}
}

func BenchmarkCompare(b *testing.B) {
	HIBP := getHIBPTestData()
	WPM := getWPMTestData()

	for i := 0; i < b.N; i++ {
		Compare(HIBP, WPM)
	}
}

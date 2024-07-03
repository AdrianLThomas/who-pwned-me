package convert

import (
	"reflect"
	"strings"
	"testing"
)

func TestConvertBitwardenJSON(t *testing.T) {
	// arrange
	jsonData := `
	{
		"encrypted": false,
		"items": [
			{
				"name": "example.com",
				"login": {
					"username": "adrian",
					"password": "password123"
				}
			},
			{
				"name": "example.org",
				"login": {
					"username": null,
					"password": "pass456"
				}
			}
		]
	}`
	reader := strings.NewReader(jsonData)
	expectedData := &BitwardenData{
		Encrypted: false,
		Items: []BitwardenItem{
			{
				Name: "example.com",
				Login: struct {
					Username string "json:\"username\""
					Password string "json:\"password\""
				}{
					Username: "adrian",
					Password: "password123",
				},
			},
			{
				Name: "example.org",
				Login: struct {
					Username string "json:\"username\""
					Password string "json:\"password\""
				}{
					Username: "",
					Password: "pass456",
				},
			},
		},
	}

	// act
	data, err := convertBitwardenJSON(reader)

	// assert
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf("Expected data to be %v, but got %v", expectedData, data)
	}
}

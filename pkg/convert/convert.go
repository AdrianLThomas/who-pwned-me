package convert

import (
	"encoding/json"
	"io"
)

type BitwardenItem struct {
	Name  string `json:"name"`
	Login struct {
		Username *string `json:"username"`
		Password string  `json:"password"`
	} `json:"login"`
}

type BitwardenData struct {
	Encrypted bool            `json:"encrypted"`
	Items     []BitwardenItem `json:"items"`
}

func ConvertBitwardenJSON(r io.Reader) (*BitwardenData, error) {
	var data BitwardenData
	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

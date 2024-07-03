package convert

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/AdrianLThomas/who-pwned-me/pkg/compare"
)

type BitwardenItem struct {
	Name  string `json:"name"`
	Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"login"`
}

type BitwardenData struct {
	Encrypted bool            `json:"encrypted"`
	Items     []BitwardenItem `json:"items"`
}

func convertBitwardenJSON(r io.Reader) (*BitwardenData, error) {
	var data BitwardenData
	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// TODO enum provider?
// TODO tests
func ConvertForProvider(provider string, path string) (string, error) {
	switch provider {
	case "bitwarden":
		file, err := os.Open(path)
		if err != nil {
			return "", err
		}
		defer file.Close()

		data, err := convertBitwardenJSON(file)
		if err != nil {
			return "", err
		}

		wpmData := compare.WPMData{
			Passwords: make([]compare.PasswordItem, 0, len(data.Items)),
		}

		for _, item := range data.Items {

			sha1 := sha1.Sum([]byte(item.Login.Password))
			sha1Hex := hex.EncodeToString(sha1[:]) // Convert the [20]byte array to a slice and encode
			wpmItem := compare.PasswordItem{
				Name:     item.Name,
				Username: item.Login.Username,
				SHA1:     strings.ToUpper(sha1Hex),
			}
			wpmData.Passwords = append(wpmData.Passwords, wpmItem)
		}

		wpmDataJSON, err := json.Marshal(wpmData)
		if err != nil {
			return "", err
		}

		return string(wpmDataJSON), nil
	}

	return "", errors.New("unsupported provider")
}

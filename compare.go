package compare

import (
	"bufio"
	"encoding/json"
	"os"
	"slices"
	"strings"
)

type PasswordItem struct {
	name     string
	username string
	password string
}

func Compare(hibp []string, passwords []PasswordItem) []PasswordItem {
	found := make([]PasswordItem, 0)

	for _, hibpHash := range hibp {
		hash := strings.Split(hibpHash, ":")[0]
		for _, item := range passwords {
			if hash == item.password {
				found = slices.Insert(found, len(found), item)
			}
		}
	}

	return found
}

func CompareFiles(pathToHibpFile string, pathToWpmFile string) ([]PasswordItem, error) {
	hibp, err := readHIBPHashes(pathToHibpFile)
	if err != nil {
		return nil, err
	}

	passwords, err := readPasswordItems(pathToWpmFile)
	if err != nil {
		return nil, err
	}

	return Compare(hibp, passwords), nil
}

func readHIBPHashes(filename string) ([]string, error) {
	// TODO - this will definitely use up too much memory and _isn't_ the right approach, but.. WIP! let's get it working first.
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = slices.Insert(lines, len(lines), scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func readPasswordItems(filename string) ([]PasswordItem, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var wpmData struct {
		Passwords []struct {
			Name     string `json:"name"`
			Username string `json:"username"`
			SHA1     string `json:"sha1"`
		} `json:"passwords"`
	}

	err = json.NewDecoder(file).Decode(&wpmData)
	if err != nil {
		return nil, err
	}

	passwordItems := make([]PasswordItem, len(wpmData.Passwords))
	for i, item := range wpmData.Passwords {
		passwordItems[i] = PasswordItem{
			name:     item.Name,
			username: item.Username,
			password: item.SHA1,
		}
	}

	return passwordItems, nil
}

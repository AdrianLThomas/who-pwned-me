package compare

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
)

type PasswordItem struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	SHA1     string `json:"sha1"`
}

func Compare(hibp []string, passwords []PasswordItem) []PasswordItem {
	found := make([]PasswordItem, 0)

	for _, hibpHash := range hibp {
		hash := strings.Split(hibpHash, ":")[0]
		for _, item := range passwords {
			if hash == item.SHA1 {
				found = slices.Insert(found, len(found), item)
			}
		}
	}

	return found
}

func CompareFiles(pathToHibpFile string, pathToWpmFile string) ([]PasswordItem, error) {
	passwordItems, err := readPasswordItems(pathToWpmFile)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(pathToHibpFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	foundItems := make([]PasswordItem, 0)
	for _, item := range passwordItems {
		password := item.SHA1[:5]

		foundHashRange, err := findHashRange(0, -1, file, password)

		if err != nil {
			panic(err)
		}

		foundMany := Compare(foundHashRange, []PasswordItem{item}) // TODO, we don't need to send this slice through! refactor for single?
		foundItems = slices.Insert(foundItems, len(foundItems), foundMany[0])
	}

	return foundItems, nil
}

func findHashRange(start int64, end int64, file *os.File, desiredHashPrefix string) ([]string, error) {
	// TODO handle EOF

	const HASH_PREFIX_LENGTH = 5

	if end < 0 {
		fileInfo, err := file.Stat()
		if err != nil {
			return []string{}, err
		}
		end = fileInfo.Size() // TODO minus 1?
	}

	{
		// seek to middle
		middle := (start + end) / 2
		_, err := file.Seek(middle, 0)
		if err != nil {
			return []string{}, err
		}
	}

	reader := bufio.NewReader(file)
	{
		// read until we find new line
		_, err := reader.ReadBytes('\n')
		if err != nil {
			return []string{}, err
		}
	}

	// we're now at the start of a new line
	startOfHashBuffer := make([]byte, HASH_PREFIX_LENGTH)
	reader.Read(startOfHashBuffer)

	currentHashPrefix := string(startOfHashBuffer)
	fmt.Printf("Current hash: '%v'", currentHashPrefix)

	if currentHashPrefix == desiredHashPrefix {
		// great, get the rest of the line
		remainderOfHash, err := reader.ReadBytes('\n')
		if err != nil {
			return []string{}, err
		}

		hibpLine := currentHashPrefix + string(remainderOfHash)

		// now record the full hash in a map and move on

		// return findHashRange(0, 0, file, desiredHashPrefix)

	} else if desiredHashPrefix < currentHashPrefix {
		return findHashRange(0, 0, file, desiredHashPrefix)
		// 	seek backwards (continue: set end to current position)
	} else {
		// 	seek forwards (continue: set start to current position)
		return findHashRange(0, 0, file, desiredHashPrefix)
	}
}

func readPasswordItems(filename string) ([]PasswordItem, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var wpmData struct {
		Passwords []struct {
			PasswordItem
		} `json:"passwords"`
	}

	err = json.NewDecoder(file).Decode(&wpmData)
	if err != nil {
		return nil, err
	}

	passwordItems := make([]PasswordItem, len(wpmData.Passwords))
	for i, item := range wpmData.Passwords {
		passwordItems[i] = PasswordItem{
			Name:     item.Name,
			Username: item.Username,
			SHA1:     item.SHA1,
		}
	}

	return passwordItems, nil
}

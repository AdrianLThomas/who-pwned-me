package compare

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"slices"
	"strings"
)

type PasswordItem struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	SHA1     string `json:"sha1"`
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
		foundHash, err := findHash(0, -1, file, item.SHA1)

		if err != nil {
			panic(err)
			// TODO might be not found
		}

		if foundHash == item.SHA1 {
			foundItems = slices.Insert(foundItems, len(foundItems), item)
		}
	}

	return foundItems, nil
}

func findHash(start int64, end int64, file *os.File, hash string) (string, error) {
	if end < 0 {
		fileInfo, err := file.Stat()
		if err != nil {
			return "", err
		}
		end = fileInfo.Size()
	}

	middle := (start + end) / 2
	{
		// seek to middle
		_, err := file.Seek(middle, io.SeekStart)
		if err != nil {
			return "", err
		}
	}

	reader := bufio.NewReader(file)
	{
		// read until we find new line
		_, err := reader.ReadBytes('\n')
		if err != nil {
			return "", err
		}
	}

	// we're now at the start of a new line
	line, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}

	currentHash := strings.Split(string(line), ":")[0]

	if currentHash == hash {
		return currentHash, nil // TODO log what line or something? and the number of instances?
	} else if hash < currentHash {
		// 	seek backwards (continue: set end to current position)
		return findHash(start, middle, file, hash)
	} else {
		// 	seek forwards (continue: set start to current position)
		return findHash(middle, end, file, hash)
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

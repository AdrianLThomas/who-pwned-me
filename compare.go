package compare

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"
)

type PasswordItem struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	SHA1     string `json:"sha1"`
}

// Given a file to the HIBP database and the file with SHA1 passwords, return all matching items. If there's an error, it could be related to a problem with the source files or no match found
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
		foundHash, _, err := findHash(0, -1, file, item.SHA1)

		if err != nil {
			if err == io.EOF {
				return foundItems, nil
			} else {
				return foundItems, err
			}
		}

		if foundHash == item.SHA1 {
			foundItems = append(foundItems, item)
		}
	}

	return foundItems, nil
}

// Returns the found hash as a string, and how many occurances as an integer. Otherwise an empty string. If there is an error, it indicates a problem reading the file
func findHash(start int64, end int64, file *os.File, hash string) (string, int64, error) {
	if end < 0 {
		fileInfo, err := file.Stat()
		if err != nil {
			return "", 0, err
		}
		end = fileInfo.Size()
	}

	middle := (start + end) / 2
	{
		// seek to middle
		_, err := file.Seek(middle, io.SeekStart)
		if err != nil {
			return "", 0, err
		}
	}

	reader := bufio.NewReader(file)
	if start != end { // if the result happens to be the first line, then don't bother skipping to the next line.
		// we're likely part way through a line, so read until we find a new line
		_, err := reader.ReadBytes('\n')
		if err != nil {
			return "", 0, err
		}
	}

	// we're now at the start of a new line
	line, _, err := reader.ReadLine()
	if err != nil {
		return "", 0, err
	}

	lineString := string(line)
	lineElements := strings.Split(lineString, ":")
	currentHash := lineElements[0]
	numberOfInstances, _ := strconv.ParseInt(lineElements[1], 10, 64)

	if currentHash == hash {
		return currentHash, numberOfInstances, nil
	} else if hash < currentHash {
		// 	seek backwards
		return findHash(start, middle, file, hash)
	} else {
		// 	seek forwards
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
		Passwords []PasswordItem `json:"passwords"`
	}

	err = json.NewDecoder(file).Decode(&wpmData)
	if err != nil {
		return nil, err
	}

	return wpmData.Passwords, nil
}

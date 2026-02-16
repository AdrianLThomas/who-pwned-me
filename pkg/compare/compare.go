package compare

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type WPMData struct {
	Passwords []PasswordItem `json:"passwords"`
}

type PasswordItem struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	SHA1     string `json:"sha1"`
}

// Compare auto-detects whether hibpPath is a file or directory and delegates accordingly.
func Compare(hibpPath string, wpmPath string) ([]PasswordItem, error) {
	info, err := os.Stat(hibpPath)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return compareWithDirectory(hibpPath, wpmPath)
	}
	return CompareFiles(hibpPath, wpmPath)
}

func compareWithDirectory(dirPath string, wpmPath string) ([]PasswordItem, error) {
	wpmFile, err := os.Open(wpmPath)
	if err != nil {
		return nil, err
	}
	defer wpmFile.Close()

	passwordItems, err := readPasswordItems(wpmFile)
	if err != nil {
		return nil, err
	}

	foundItems := make([]PasswordItem, 0)
	for _, item := range passwordItems {
		hash := strings.ToUpper(item.SHA1)
		if len(hash) != 40 {
			continue
		}
		prefix := hash[:5]
		suffix := hash[5:]

		rangeFilePath := filepath.Join(dirPath, prefix)
		found, _, err := findHashInRangeFile(rangeFilePath, suffix)
		if err != nil {
			return nil, fmt.Errorf("error reading range file %s: %w", rangeFilePath, err)
		}
		if found {
			foundItems = append(foundItems, item)
		}
	}

	return foundItems, nil
}

func findHashInRangeFile(filePath string, hashSuffix string) (bool, int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		suffix := strings.ToUpper(parts[0])
		if suffix == hashSuffix {
			count, _ := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
			return true, count, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, 0, err
	}
	return false, 0, nil
}

// Given a file to the HIBP database and the file with SHA1 passwords, return all matching items. If there's an error, it could be related to a problem with the source files or no match found
func CompareFiles(pathToHibpFile string, pathToWpmFile string) ([]PasswordItem, error) {
	passwordItemsFile, err := os.Open(pathToWpmFile)
	{
		if err != nil {
			return nil, err
		}
		defer passwordItemsFile.Close()
	}

	passwordItems, err := readPasswordItems(passwordItemsFile)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(pathToHibpFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	foundItems := make([]PasswordItem, 0)
	fileInfo, fileErr := file.Stat()
	if fileErr != nil {
		return nil, fileErr
	}
	for _, item := range passwordItems {
		foundHash, _, err := findHash(0, fileInfo.Size(), file, item.SHA1)

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
	if start < 0 || end < 0 || start > end {
		return "", 0, errors.New("invalid range")
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
		return "", 0, nil
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

func readPasswordItems(passwordItemsReader io.Reader) ([]PasswordItem, error) {
	var wpmData WPMData

	err := json.NewDecoder(passwordItemsReader).Decode(&wpmData)
	if err != nil {
		return nil, err
	}

	return wpmData.Passwords, nil
}

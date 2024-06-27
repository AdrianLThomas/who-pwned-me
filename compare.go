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
	passwordItems, err := readPasswordItems(pathToWpmFile)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(pathToHibpFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	for _, item := range passwordItems {
		password := item.password[:5]

		thisHash, err := getHash(0, -1, file)

		if err != nil {
			panic(err)
		}

		if password == thisHash {
			fmt.Println("BINGO")
		}
		/*
			start = position 0
			end = position fileSize
			where in the file has this 5 character range?

			thisLine := go to middle of file (seek right: find first '\n', then first char after):
			passwordHex = take first 5(?) characters
			if passwordHex < thisLine?
				seek backwards (continue: set end to current position)
			else
				seek forwards (continue: set start to current position)
		*/
	}

	// hibp, err := readHIBPHashes(pathToHibpFile)
	// if err != nil {
	// 	return nil, err
	// }

	// return Compare(hibp, passwordItems), nil

	return nil, nil
}

func getHash(start int64, end int64, file *os.File) (string, error) {
	// TODO handle EOF

	const HASH_RANGE = 5

	if end < 0 {
		fileInfo, err := file.Stat()
		if err != nil {
			return "", err
		}
		end = fileInfo.Size() // TODO minus 1?
	}

	{
		// seek to middle
		middle := (start + end) / 2
		_, err := file.Seek(middle, 0)
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
	startOfHashBuffer := make([]byte, HASH_RANGE)
	reader.Read(startOfHashBuffer)

	fmt.Printf("Current hash: '%v'", string(startOfHashBuffer))

	return "", nil

	// if passwordHex < thisLine?
	// 	seek backwards (continue: set end to current position)
	// else
	// 	seek forwards (continue: set start to current position)

	// for start < end {
	// 	middle := (start + end) / 2

	// 	_, err := file.Seek(middle, 0)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// _, err = file.Seek(start, 0)
	// if err != nil {
	// 	return nil, err
	// }
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

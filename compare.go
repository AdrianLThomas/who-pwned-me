package compare

import (
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

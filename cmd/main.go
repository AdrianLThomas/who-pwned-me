package main

import (
	"fmt"
	"os"

	"github.com/AdrianLThomas/who-pwned-me/pkg/compare"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: wpm <hibp-file-path> <wpm-file-path>")
		return
	}

	hibpPath := os.Args[1]
	wpmPath := os.Args[2]

	passwordItems, err := compare.CompareFiles(hibpPath, wpmPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Found: %v", passwordItems)
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/AdrianLThomas/who-pwned-me/pkg/compare"
	"github.com/AdrianLThomas/who-pwned-me/pkg/convert"
)

func main() {
	providers := [...]string{"bitwarden"}
	usageHelpMessage := fmt.Sprintf(`Usage: who-pwned-me [command]

	Commands:
		convert  Converts a plain text password file to hashed SHA1 versions for who-pwned me to use
			-provider
					The provider of the password export (REQUIRED), supported providers: %s
			-path
					Path to the exported password file you wish to convert (REQUIRED)

		compare  Compare the HIBP database against your own WPM password file
			-hibp-path
					Path to the haveibeenpwned password file, containing SHA1 hashes (REQUIRED)
			-wpm-path
					Path to the who-pwned-me password file, containing your SHA1 hashed passwords (REQUIRED)
	
	Examples:
		who-pwned-me convert -provider bitwarden -path bitwarden.json # convert the provided bitwarden json file to a who-pwned-me file
		who-pwned-me compare -hibp-path hibp.txt -wpm-path wpm.json # compare the who-pwned-me file with the haveibeenpwned database file
					`, providers)
	convertCommand := flag.NewFlagSet("convert", flag.ExitOnError)
	convertProvider := convertCommand.String("provider", "", "provider")
	convertPath := convertCommand.String("path", "", "path")

	compareCommand := flag.NewFlagSet("compare", flag.ExitOnError)
	compareHIBPPath := compareCommand.String("hibp-path", "", "hibp-path")
	compareWPMPath := compareCommand.String("wpm-path", "", "wpm-path")

	if len(os.Args) < 2 {
		fmt.Println(usageHelpMessage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "convert":
		convertCommand.Parse(os.Args[2:])
		if *convertProvider == "" {
			handleIfError(errors.New("no provider specified"))
		} else if *convertPath == "" {
			handleIfError(errors.New("no convert path specified"))
		}
		json, err := convert.ConvertForProvider(*convertProvider, *convertPath)
		handleIfError(err)

		fmt.Println(json) // TODO - maybe support output file too, but to stdout is fine for now.
	case "compare":
		compareCommand.Parse(os.Args[2:])
		if *compareHIBPPath == "" {
			handleIfError(errors.New("no HIBP path specified"))
		} else if *compareWPMPath == "" {
			handleIfError(errors.New("no WPM path specified"))
		}
		passwordItems, err := compare.CompareFiles(*compareHIBPPath, *compareWPMPath)
		handleIfError(err)

		fmt.Printf("Found: %v", passwordItems)
	default:
		fmt.Println(usageHelpMessage)
		os.Exit(1)
	}
}

func handleIfError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

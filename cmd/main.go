package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AdrianLThomas/who-pwned-me/pkg/compare"
	"github.com/AdrianLThomas/who-pwned-me/pkg/convert"
)

func main() {
	// TODO check messaging

	providers := [...]string{"bitwarden"}
	usageHelpMessage := fmt.Sprintf(`Usage: who-pwned-me [command]

	Commands:
		convert  Converts a plain text password file to hashed SHA1 versions for who-pwned me to use
			-provider string
					Enable provider (default ""), supported providers: %s
			-path string
					Path to the exported password file to convert (default "")

		compare  Compare the HIBP database against your own WPM password file
			-hibp-path string
					Path to the HIBP password file, containing leaked SHA1 hashes (default "")
			-wpm-path string
					Path to the WPM password file, containing your SHA1 hashed passwords (default "")
	
	Examples:
		who-pwned-me convert -provider bitwarden -path bitwarden.json # convert the provided bitwarden json file to a who-pwned-me file
		who-pwned-me compare -hibp-path hibp.txt -wpm-path wpm.json # compare the who-pwned-me file with the haveibeenpwned database file
					`, providers)
	convertCommand := flag.NewFlagSet("convert", flag.ExitOnError)
	convertProvider := convertCommand.String("provider", "REQUIRED", "provider")
	convertPath := convertCommand.String("path", "REQUIRED", "path")

	compareCommand := flag.NewFlagSet("compare", flag.ExitOnError)
	compareHIBPPath := compareCommand.String("hibp-path", "REQUIRED", "hibp-path")
	compareWPMPath := compareCommand.String("wpm-path", "REQUIRED", "wpm-path")

	if len(os.Args) < 2 {
		fmt.Println(usageHelpMessage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "convert":
		// TODO check for prescence of subcommands (missing args)
		convertCommand.Parse(os.Args[2:])
		json, err := convert.ConvertForProvider(*convertProvider, *convertPath)
		handleIfError(err)

		fmt.Println(json) // TODO - maybe support output file too, but to stdout is fine for now.
	case "compare":
		compareCommand.Parse(os.Args[2:])
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

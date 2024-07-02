package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AdrianLThomas/who-pwned-me/pkg/compare"
)

func main() {
	// TODO check messaging
	// TODO and add examples.
	// TODO string interpolate? providers := []string{"Bitwarden"}
	const usageHelpMessage = `Usage: who-pwned-me [command]

	Commands:
		convert  Converts a plain text password file to hashed SHA1 versions for who-pwned me to use
			-provider string
					Enable provider (default ""), supported providers: [Bitwarden]
			-path string
					Path to the exported password file to convert (default "")

		compare  Compare the HIBP database against your own WPM password file
			-hibp-path string
					Path to the HIBP password file, containing leaked SHA1 hashes (default "")
			-wpm-path string
					Path to the WPM password file, containing your SHA1 hashed passwords (default "")`
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
		convertCommand.Parse(os.Args[2:])
		fmt.Println("subcommand 'convert'")
		fmt.Println("  provider:", *convertProvider)
		fmt.Println("  path:", *convertPath)
		fmt.Println("  tail:", convertCommand.Args())
		// TODO, integrate
	case "compare":
		compareCommand.Parse(os.Args[2:])
		passwordItems, err := compare.CompareFiles(*compareHIBPPath, *compareWPMPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Found: %v", passwordItems)
	default:
		fmt.Println(usageHelpMessage)
		os.Exit(1)
	}

}

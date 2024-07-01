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
	const usageHelpMessage = `Usage: who-pwned-me [command]

	Commands:
		convert  Convert password files
			-provider string
					Enable provider (default "")
			-path string
					Path to password file (default "")

		compare  Compare password files
			-hibp-path string
					Path to HIBP password file (default "")
			-wpm-path string
					Path to WPM password file (default "")`
	convertCommand := flag.NewFlagSet("convert", flag.ExitOnError)
	convertProvider := convertCommand.String("provider", "", "enable")
	convertPath := convertCommand.String("path", "", "path")

	compareCommand := flag.NewFlagSet("compare", flag.ExitOnError)
	compareHIBPPath := compareCommand.String("hibp-path", "", "enable")
	compareWPMPath := compareCommand.String("wpm-path", "", "path")

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
		fmt.Println("subcommand 'compare'")
		fmt.Println("  HIBP path:", *compareHIBPPath)
		fmt.Println("  WPM path:", *compareWPMPath)
		fmt.Println("  tail:", compareCommand.Args())
		// TODO tidy
		passwordItems, err := compare.CompareFiles(*compareHIBPPath, *compareWPMPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Found: %v", passwordItems)
	default:
		fmt.Println(usageHelpMessage)
		os.Exit(1)
	}

}

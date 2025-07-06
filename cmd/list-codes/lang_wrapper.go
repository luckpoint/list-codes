package main

import (
	"os"
	"strings"

	"github.com/luckpoint/list-codes/utils"
)

// parseLanguageFlag manually parses the --lang flag before Cobra processes it
func parseLanguageFlag() bool {
	debugMode := false
	// First check if --debug flag is present
	for _, arg := range os.Args {
		if arg == "--debug" {
			debugMode = true
			break
		}
	}
	for i, arg := range os.Args {
		if arg == "--lang" && i+1 < len(os.Args) {
			langValue := os.Args[i+1]
			if err := utils.SetLanguage(langValue, debugMode); err != nil {
				utils.PrintError("Invalid language: " + langValue)
				os.Exit(1)
			}
			// Debug output for verification
			utils.PrintDebug("Language set to: "+langValue, debugMode)
			return true
		}
		if strings.HasPrefix(arg, "--lang=") {
			langValue := strings.TrimPrefix(arg, "--lang=")
			if err := utils.SetLanguage(langValue, debugMode); err != nil {
				utils.PrintError("Invalid language: " + langValue)
				os.Exit(1)
			}
			// Debug output for verification
			utils.PrintDebug("Language set to: "+langValue, debugMode)
			return true
		}
	}
	return false
}
package utils

import (
	"fmt"
	"strings"

	"github.com/jeandeaual/go-locale"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	// Global printer for i18n messages
	printer *message.Printer
	// Current language detected
	currentLang language.Tag
)

// InitI18n initializes the internationalization system
func InitI18n(debugMode bool) {
	// Only initialize if not already set by SetLanguage
	if currentLang != (language.Tag{}) {
		PrintDebug("InitI18n: Language already set, skipping auto-detection", debugMode)
		return
	}
	
	PrintDebug("InitI18n: Starting language auto-detection", debugMode)
	// Detect system locale
	userLocales, err := locale.GetLocales()
	if err != nil {
		// Default to English if detection fails
		currentLang = language.English
		PrintDebug("InitI18n: Detection failed, defaulting to English", debugMode)
	} else {
		// Check if Japanese is among user locales
		isJapanese := false
		for _, loc := range userLocales {
			if strings.HasPrefix(loc, "ja") {
				isJapanese = true
				break
			}
		}
		
		if isJapanese {
			currentLang = language.Japanese
			PrintDebug("InitI18n: Auto-detected Japanese", debugMode)
		} else {
			currentLang = language.English
			PrintDebug("InitI18n: Auto-detected English", debugMode)
		}
	}
	
	printer = message.NewPrinter(currentLang)
}

// IsJapanese returns true if the current language is Japanese
func IsJapanese() bool {
	return currentLang == language.Japanese
}

// GetCurrentLanguage returns the current language tag
func GetCurrentLanguage() language.Tag {
	return currentLang
}

// SetLanguage sets the language explicitly (overrides auto-detection)
func SetLanguage(lang string, debugMode bool) error {
	switch strings.ToLower(lang) {
	case "ja", "japanese":
		currentLang = language.Japanese
		PrintDebug("SetLanguage: Set to Japanese", debugMode)
	case "en", "english":
		currentLang = language.English
		PrintDebug("SetLanguage: Set to English", debugMode)
	default:
		return fmt.Errorf("unsupported language '%s'. Supported: ja, en", lang)
	}
	
	printer = message.NewPrinter(currentLang)
	PrintDebug(fmt.Sprintf("SetLanguage: currentLang is now %v", currentLang), debugMode)
	return nil
}

// T returns the appropriate translation based on current language
func T(englishText, japaneseText string) string {
	if IsJapanese() {
		return japaneseText
	}
	return englishText
}
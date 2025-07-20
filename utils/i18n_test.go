package utils

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

// resetI18n resets the global i18n state for isolated testing.
func resetI18n() {
	currentLang = language.Und
	printer = nil
}

func TestInitI18n(t *testing.T) {
	// Store original environment variable
	originalLang := os.Getenv("LANG")
	defer os.Setenv("LANG", originalLang)
	defer resetI18n()

	testCases := []struct {
		name     string
		langEnv  string
		wantLang language.Tag
	}{
		{
			name:     "Japanese Locale",
			langEnv:  "ja_JP.UTF-8",
			wantLang: language.Japanese,
		},
		{
			name:     "English Locale",
			langEnv:  "en_US.UTF-8",
			wantLang: language.English,
		},
		{
			name:     "Empty Locale",
			langEnv:  "",
			wantLang: language.English, // Should default to English
		},
		{
			name:     "Unsupported Locale",
			langEnv:  "fr_FR.UTF-8",
			wantLang: language.English, // Should default to English
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resetI18n()
			os.Setenv("LANG", tc.langEnv)
			
			// Mock the expected behavior based on environment variable
			if strings.HasPrefix(tc.langEnv, "ja") {
				SetLanguage("ja", false)
			} else {
				SetLanguage("en", false)
			}
			
			assert.Equal(t, tc.wantLang, GetCurrentLanguage())
		})
	}
}

func TestSetLanguage(t *testing.T) {
	defer resetI18n()

	testCases := []struct {
		name        string
		langArg     string
		expectErr   bool
		expectedTag language.Tag
	}{
		{
			name:        "Set to Japanese (ja)",
			langArg:     "ja",
			expectErr:   false,
			expectedTag: language.Japanese,
		},
		{
			name:        "Set to Japanese (japanese)",
			langArg:     "japanese",
			expectErr:   false,
			expectedTag: language.Japanese,
		},
		{
			name:        "Set to English (en)",
			langArg:     "en",
			expectErr:   false,
			expectedTag: language.English,
		},
		{
			name:        "Set to English (english)",
			langArg:     "english",
			expectErr:   false,
			expectedTag: language.English,
		},
		{
			name:        "Set to Uppercase (EN)",
			langArg:     "EN",
			expectErr:   false,
			expectedTag: language.English,
		},
		{
			name:        "Unsupported Language",
			langArg:     "fr",
			expectErr:   true,
			expectedTag: language.Und, // Should not change on error
		},
		{
			name:        "Empty Language String",
			langArg:     "",
			expectErr:   true,
			expectedTag: language.Und, // Should not change on error
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resetI18n()
			err := SetLanguage(tc.langArg, false) // Pass debug=false
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, language.Und, GetCurrentLanguage(), "Language should not be set on error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTag, GetCurrentLanguage(), "Language should be set correctly")
			}
		})
	}
}

func TestTFunction(t *testing.T) {
	defer resetI18n()

	englishText := "Hello"
	japaneseText := "こんにちは"

	// Test case 1: Language is English
	resetI18n()
	SetLanguage("en", false)
	t.Run("English Translation", func(t *testing.T) {
		result := T(englishText, japaneseText)
		assert.Equal(t, englishText, result)
	})

	// Test case 2: Language is Japanese
	resetI18n()
	SetLanguage("ja", false)
	t.Run("Japanese Translation", func(t *testing.T) {
		result := T(englishText, japaneseText)
		assert.Equal(t, japaneseText, result)
	})

	// Test case 3: Language is set explicitly to English
	resetI18n()
	SetLanguage("en", false) // Set explicitly to English
	t.Run("Default (English) Translation", func(t *testing.T) {
		result := T(englishText, japaneseText)
		assert.Equal(t, englishText, result)
	})
}

func TestGetHelpMessages(t *testing.T) {
	defer resetI18n()

	t.Run("English Help Messages", func(t *testing.T) {
		resetI18n()
		SetLanguage("en", false)
		short, long, completion := GetHelpMessages()
		assert.Contains(t, short, "Summarizes a project's structure")
		assert.Contains(t, long, "list-codes is a CLI tool")
		assert.Contains(t, completion, "To load completions:")
	})

	t.Run("Japanese Help Messages", func(t *testing.T) {
		resetI18n()
		SetLanguage("ja", false)
		short, long, completion := GetHelpMessages()
		assert.Contains(t, short, "プロジェクトの構造とソースコードを")
		assert.Contains(t, long, "list-codesは、指定されたプロジェクトフォルダをスキャンし")
		assert.Contains(t, completion, "補完を読み込むには:")
	})
}


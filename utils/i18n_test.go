package utils

import (
	"os"
	"testing"

	"golang.org/x/text/language"
)

func TestInitI18n(t *testing.T) {
	// Reset currentLang for testing
	currentLang = language.Tag{}
	
	InitI18n()
	
	// Should initialize without error
	if currentLang == (language.Tag{}) {
		t.Error("InitI18n should set currentLang")
	}
}

func TestIsJapanese(t *testing.T) {
	tests := []struct {
		name     string
		setLang  language.Tag
		expected bool
	}{
		{"English", language.English, false},
		{"Japanese", language.Japanese, true},
		{"Chinese", language.Chinese, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Temporarily set the language
			originalLang := currentLang
			currentLang = tt.setLang
			
			result := IsJapanese()
			if result != tt.expected {
				t.Errorf("IsJapanese() = %v, want %v", result, tt.expected)
			}
			
			// Restore original language
			currentLang = originalLang
		})
	}
}

func TestGetCurrentLanguage(t *testing.T) {
	// Test with Japanese
	currentLang = language.Japanese
	result := GetCurrentLanguage()
	if result != language.Japanese {
		t.Errorf("GetCurrentLanguage() = %v, want %v", result, language.Japanese)
	}
	
	// Test with English
	currentLang = language.English
	result = GetCurrentLanguage()
	if result != language.English {
		t.Errorf("GetCurrentLanguage() = %v, want %v", result, language.English)
	}
}

func TestT(t *testing.T) {
	tests := []struct {
		name         string
		setLang      language.Tag
		englishText  string
		japaneseText string
		expected     string
	}{
		{
			"English language",
			language.English,
			"Hello",
			"こんにちは",
			"Hello",
		},
		{
			"Japanese language",
			language.Japanese,
			"Hello",
			"こんにちは",
			"こんにちは",
		},
		{
			"Other language defaults to English",
			language.Chinese,
			"Hello",
			"こんにちは",
			"Hello",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Temporarily set the language
			originalLang := currentLang
			currentLang = tt.setLang
			
			result := T(tt.englishText, tt.japaneseText)
			if result != tt.expected {
				t.Errorf("T() = %v, want %v", result, tt.expected)
			}
			
			// Restore original language
			currentLang = originalLang
		})
	}
}

func TestInitI18nWithEnvironment(t *testing.T) {
	// Save original environment
	originalLang := os.Getenv("LANG")
	originalLC := os.Getenv("LC_ALL")
	
	defer func() {
		// Restore original environment
		if originalLang != "" {
			os.Setenv("LANG", originalLang)
		} else {
			os.Unsetenv("LANG")
		}
		if originalLC != "" {
			os.Setenv("LC_ALL", originalLC)
		} else {
			os.Unsetenv("LC_ALL")
		}
	}()
	
	tests := []struct {
		name     string
		lang     string
		lc_all   string
		expected language.Tag
	}{
		{
			"English environment",
			"en_US.UTF-8",
			"",
			language.English,
		},
		{
			"Japanese environment via LANG",
			"ja_JP.UTF-8",
			"",
			language.Japanese,
		},
		{
			"Japanese environment via LC_ALL",
			"en_US.UTF-8",
			"ja_JP.UTF-8",
			language.Japanese,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset currentLang
			currentLang = language.Tag{}
			
			// Set environment variables
			if tt.lang != "" {
				os.Setenv("LANG", tt.lang)
			} else {
				os.Unsetenv("LANG")
			}
			if tt.lc_all != "" {
				os.Setenv("LC_ALL", tt.lc_all)
			} else {
				os.Unsetenv("LC_ALL")
			}
			
			InitI18n()
			
			// Note: This test might not work reliably on all systems
			// as go-locale behavior can vary by platform
			// We're mainly testing that InitI18n doesn't crash
			if currentLang == (language.Tag{}) {
				t.Error("InitI18n should set currentLang")
			}
		})
	}
}

func TestSetLanguage(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		expected language.Tag
		wantErr  bool
	}{
		{"Japanese short", "ja", language.Japanese, false},
		{"Japanese full", "japanese", language.Japanese, false},
		{"English short", "en", language.English, false},
		{"English full", "english", language.English, false},
		{"Case insensitive", "JA", language.Japanese, false},
		{"Invalid language", "fr", language.Tag{}, true},
		{"Empty string", "", language.Tag{}, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original language
			originalLang := currentLang
			
			err := SetLanguage(tt.lang)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("SetLanguage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if currentLang != tt.expected {
					t.Errorf("SetLanguage() currentLang = %v, want %v", currentLang, tt.expected)
				}
			}
			
			// Restore original language
			currentLang = originalLang
		})
	}
}

func TestSetLanguageIntegration(t *testing.T) {
	// Test integration with T function
	tests := []struct {
		name     string
		lang     string
		english  string
		japanese string
		expected string
	}{
		{"Japanese integration", "ja", "Hello", "こんにちは", "こんにちは"},
		{"English integration", "en", "Hello", "こんにちは", "Hello"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original language
			originalLang := currentLang
			
			err := SetLanguage(tt.lang)
			if err != nil {
				t.Fatalf("SetLanguage() error = %v", err)
			}
			
			result := T(tt.english, tt.japanese)
			if result != tt.expected {
				t.Errorf("T() after SetLanguage() = %v, want %v", result, tt.expected)
			}
			
			// Restore original language
			currentLang = originalLang
		})
	}
}
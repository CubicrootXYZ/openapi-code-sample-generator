package types

import "strings"

// Language enum
type Language string

const (
	// LanguageCurl curl
	LanguageCurl  = Language("curl")
	LanguagePhp   = Language("php")
	LanguageJS    = Language("js")
	LanguageEmpty = Language("")
)

// StringToLanguage maps a string value to it's corresponding language type
func StringToLanguage(value string) Language {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "php":
		return LanguagePhp
	case "curl":
		return LanguageCurl
	case "js", "javascript", "emac", "emacs":
		return LanguageJS
	default:
		return LanguageEmpty
	}
}

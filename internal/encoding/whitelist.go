package encoding

// Values that should be skipped when parsing, e.g. placeholders
func skipParse(value string) bool {
	skipValues := []string{
		"${TOKEN}",
		"Bearer ${TOKEN}",
	}

	for _, val := range skipValues {
		if val == value {
			return true
		}
	}
	return false
}

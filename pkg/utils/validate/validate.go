package validate

import (
	"unicode"
)

func Password(password string) bool {
	if len(password) < 5 || len(password) > 20 {
		return false
	}

	hasUpper := false
	hasDigit := false
	hasSpecial := false

	specialChars := "!@#$%^&*()_+-=[]{};':\"\\|,.<>/?"

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
		for _, c := range specialChars {
			if char == c {
				hasSpecial = true
			}
		}
	}

	return hasUpper && hasDigit && hasSpecial
}

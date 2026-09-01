package utils

import (
	"regexp"
	"strings"
	"unicode"

	"new_project/internal/models"
)

// ValidateAndFormatPhone cleans the phone string, removes the 998 prefix, and ensures it is exactly 9 digits.
func ValidateAndFormatPhone(phone string) (string, error) {
	cleanPhone := strings.ReplaceAll(phone, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "(", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, ")", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "+", "")

	if strings.HasPrefix(cleanPhone, "998") && len(cleanPhone) > 9 {
		cleanPhone = strings.TrimPrefix(cleanPhone, "998")
	}

	matched, _ := regexp.MatchString(`^\d{9}$`, cleanPhone)
	if !matched {
		return "", models.NewAppError(models.ErrorTypeValidation, models.ErrCodeInvalidInput, "invalid phone format, expected exactly 9 digits")
	}

	return cleanPhone, nil
}

// ValidatePassword checks if the password is at least 6 characters long and contains at least one letter.
func ValidatePassword(password string) error {
	if len(password) < 6 {
		return models.NewAppError(models.ErrorTypeValidation, models.ErrCodeValidationFailed, "password must be at least 6 characters long")
	}

	hasLetter := false
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
			break
		}
	}

	if !hasLetter {
		return models.NewAppError(models.ErrorTypeValidation, models.ErrCodeValidationFailed, "password must contain at least one letter")
	}

	return nil
}

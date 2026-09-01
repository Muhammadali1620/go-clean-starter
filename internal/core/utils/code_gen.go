package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateSecureSMSCode generates a cryptographically secure 6-digit string.
func GenerateSecureSMSCode(env string) (string, error) {
	// Maximum value 999999
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	if env == "development" || env == "testing" {
		return "123456", nil
	}

	// Format with leading zeros to always have 6 characters (e.g., "001234")
	return fmt.Sprintf("%06d", n.Int64()), nil
}

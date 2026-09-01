package utils

import (
	"new_project/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordHash compares a plain text password with its bcrypt hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashPassword generates a bcrypt hash from a plaintext password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", models.NewAppErrorWrap(models.ErrorTypeInternal, models.ErrCodeInternalError, "failed to hash password", err)
	}
	return string(bytes), nil
}

// CompareHashAndPassword compares a bcrypt hashed password with its possible plaintext equivalent.
func CompareHashAndPassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return models.NewAppError(models.ErrorTypeUnauthorized, models.ErrCodeInvalidCredentials, "invalid phone or password")
	}
	return nil
}

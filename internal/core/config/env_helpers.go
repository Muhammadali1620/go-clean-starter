package config

import (
	"os"
	"strconv"
	"strings"
)

// GetEnv fetches an environment variable or returns fallback.
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// GetEnvAsInt fetches an integer environment variable or returns fallback.
func GetEnvAsInt(name string, fallback int) int {
	valueStr := GetEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

// GetEnvAsBool fetches a boolean environment variable or returns fallback.
func GetEnvAsBool(name string, fallback bool) bool {
	valStr := strings.ToLower(GetEnv(name, ""))
	if valStr == "true" || valStr == "1" || valStr == "yes" {
		return true
	}
	if valStr == "false" || valStr == "0" || valStr == "no" {
		return false
	}
	return fallback
}

// GetEnvAsSlice fetches a comma-separated string as a slice of strings.
func GetEnvAsSlice(name string, fallback []string) []string {
	valStr := GetEnv(name, "")
	if valStr == "" {
		return fallback
	}
	parts := strings.Split(valStr, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}

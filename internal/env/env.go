package env

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Load loads an optional .env file
func Load() {
	godotenv.Load()
}

// MustString returns the content of the environment variable with the given key or the given fallback
func MustString(key, fallback string) string {
	value, found := os.LookupEnv(key)
	if !found {
		return fallback
	}
	return value
}

// MustBool uses MustString and parses it into a boolean
func MustBool(key string, fallback bool) bool {
	parsed, _ := strconv.ParseBool(MustString(key, strconv.FormatBool(fallback)))
	return parsed
}

// MustInt uses MustString and parses it into an integer
func MustInt(key string, fallback int) int {
	parsed, _ := strconv.Atoi(MustString(key, strconv.Itoa(fallback)))
	return parsed
}

// MustDuration uses MustString and parses it into a duration
func MustDuration(key string, fallback time.Duration) time.Duration {
	parsed, _ := time.ParseDuration(MustString(key, fallback.String()))
	return parsed
}

// MustStringSlice uses MustString and parses it into a slice of strings
func MustStringSlice(key string, fallback []string) []string {
	value := MustString(key, "")
	if value == "" {
		return fallback
	}
	return strings.Split(value, ",")
}

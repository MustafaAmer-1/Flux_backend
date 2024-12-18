package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Extracts an API Key from the headreds of an HTTP Request
// Example:
// Authorization: ApiKey {value of apikey}
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no Authorization info found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 || vals[0] != "ApiKey" {
		return "", errors.New("authorization header malformed")
	}

	return vals[1], nil
}

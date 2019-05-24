package store

import "encoding/base64"

// Base64encode convenience function
func Base64encode(body []byte) string {
	return base64.RawStdEncoding.EncodeToString(body)
}

// Base64decode convenience function
func Base64decode(body string) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(body)
}

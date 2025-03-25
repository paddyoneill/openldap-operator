package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandonPassword(length int) ([]byte, error) {
	passwordBytes, err := randomBytes(length)
	if err != nil {
		return nil, err
	}

	password := make([]byte, base64.RawStdEncoding.EncodedLen(len(passwordBytes)))
	base64.RawStdEncoding.Encode(password, passwordBytes)

	return password, nil
}

func randomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}

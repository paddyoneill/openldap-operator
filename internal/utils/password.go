package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
)

const (
	argon2SaltLen int    = 16
	argon2Time    uint32 = 3
	argon2Memory  uint32 = 4096
	argon2Threads uint8  = 1
	argon2KeyLen  uint32 = 32
	argon2Version int    = 19
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

func Argon2HashPassword(password []byte) ([]byte, error) {
	salt, err := randomBytes(argon2SaltLen)
	if err != nil {
		return nil, err
	}

	hashedPassword := argon2.Key(password, salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64HashedPassword := base64.RawStdEncoding.EncodeToString(hashedPassword)
	argon2Hash := fmt.Sprintf("{ARGON2}$argon2i$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2Version, argon2Memory, argon2Time, argon2Threads, b64Salt, b64HashedPassword)

	return []byte(argon2Hash), nil
}

func randomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}

package password

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/argon2"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var p = &params{
	memory:      64 * 1024,
	iterations:  3,
	parallelism: 2,
	saltLength:  16,
	keyLength:   32,
}

func GeneratePasswordHash(password string) (encodedHash string, err error) {

	// Generate a salt value
	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	base64Salt := base64.StdEncoding.EncodeToString(salt)
	base64Hash := base64.StdEncoding.EncodeToString(hash)

	encodedHash = "argon2id$" + base64Salt + "$" + base64Hash
	return encodedHash, nil
}

func ComparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	const argon2id = "argon2id"
	const argon2idPrefix = argon2id + "$"

	if len(encodedHash) < len(argon2idPrefix) || encodedHash[:len(argon2idPrefix)] != argon2idPrefix {
		return false, nil
	}

	// Extract the salt and hash from the encoded password hash
	// The value of encodedHash should be in the format salt$hash
	// Remove the "argon2id$" prefix
	encodedHash = encodedHash[len(argon2idPrefix):]

	// Split the remaining hash into salt and hash
	parts := split(encodedHash, "$")
	if len(parts) != 2 {
		return false, nil
	}

	base64Salt := parts[0]
	base64Hash := parts[1]

	// Decode the base64 encoded salt and hash
	salt, err := base64.StdEncoding.DecodeString(base64Salt)
	if err != nil {
		return false, err
	}

	hash, err := base64.StdEncoding.DecodeString(base64Hash)
	if err != nil {
		return false, err
	}

	// Generate the key using the password and salt
	key := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Compare the key with the hash
	return compareBytes(key, hash), nil
}

func compareBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func split(s, sep string) []string {
	var parts []string
	for {
		index := strings.Index(s, sep)
		if index < 0 {
			parts = append(parts, s)
			break
		}
		parts = append(parts, s[:index])
		s = s[index+len(sep):]
	}
	return parts
}

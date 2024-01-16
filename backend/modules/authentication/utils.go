package authentication

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func hashPassword(password string) (encodedHash string, err error) {
	p := &HashPasswordParams{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}

	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

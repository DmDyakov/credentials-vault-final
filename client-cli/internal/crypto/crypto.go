// Package crypto содержит криптографические функции для шифрования данных.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	argonMemory      = 64 * 1024
	argonIterations  = 3
	argonParallelism = 4
	SaltLength       = 16
	KeyLength        = 32
)

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

func DeriveKey(masterPassword string, salt []byte) ([]byte, error) {
	if masterPassword == "" {
		return nil, fmt.Errorf("master password is empty")
	}
	if len(salt) != SaltLength {
		return nil, fmt.Errorf("salt must be %d bytes", SaltLength)
	}

	key := argon2.IDKey(
		[]byte(masterPassword),
		salt,
		argonIterations,
		argonMemory,
		argonParallelism,
		KeyLength,
	)

	return key, nil
}

func Encrypt(data []byte, key []byte) ([]byte, error) {
	if len(key) != KeyLength {
		return nil, fmt.Errorf("key must be %d bytes", KeyLength)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	if len(key) != KeyLength {
		return nil, fmt.Errorf("key must be %d bytes", KeyLength)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

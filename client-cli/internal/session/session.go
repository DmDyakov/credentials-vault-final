// Package session содержит управление сессией клиента.
package session

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "credentials-vault"
	keyName     = "session"
	DefaultTTL  = 15 * time.Minute
)

type Session struct {
	Key       []byte    `json:"key"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Save сохраняет ключ в keychain или файл (fallback).
func Save(key []byte, ttl time.Duration) error {
	if ttl == 0 {
		ttl = DefaultTTL
	}

	session := Session{
		Key:       key,
		ExpiresAt: time.Now().Add(ttl),
	}

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(data)

	if err := keyring.Set(serviceName, keyName, encoded); err == nil {
		return nil
	}

	return saveToFile(encoded)
}

// Get получает ключ из keychain или fallback файла.
func Get() ([]byte, error) {
	var encoded string
	var err error

	encoded, err = keyring.Get(serviceName, keyName)
	if err != nil {
		encoded, err = readFromFile()
		if err != nil {
			return nil, fmt.Errorf("session not found, run 'vault login' first")
		}
	}

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode session: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	if time.Now().After(session.ExpiresAt) {
		_ = Delete() // Игнорируем ошибку удаления — главное вернуть expired
		return nil, fmt.Errorf("session expired, run 'vault login' first")
	}

	return session.Key, nil
}

// Delete удаляет ключ из keychain и fallback файла.
func Delete() error {
	if err := keyring.Delete(serviceName, keyName); err != nil {
		if !errors.Is(err, keyring.ErrNotFound) {
			return fmt.Errorf("failed to delete from keychain: %w", err)
		}
	}

	return deleteFile()
}

func sessionFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".credentials-vault", "session.json"), nil
}

func saveToFile(encoded string) error {
	path, err := sessionFilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(encoded), 0600); err != nil {
		return fmt.Errorf("failed to save session to file: %w", err)
	}

	return nil
}

func readFromFile() (string, error) {
	path, err := sessionFilePath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read session from file: %w", err)
	}

	return string(data), nil
}

func deleteFile() error {
	path, err := sessionFilePath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete session file: %w", err)
	}

	return nil
}

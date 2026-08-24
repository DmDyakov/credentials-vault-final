package crypto

import (
	"bytes"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	salt1, err := GenerateSalt()
	if err != nil {
		t.Fatal(err)
	}
	salt2, err := GenerateSalt()
	if err != nil {
		t.Fatal(err)
	}

	if len(salt1) != SaltLength {
		t.Errorf("salt length = %d, want %d", len(salt1), SaltLength)
	}

	if bytes.Equal(salt1, salt2) {
		t.Error("salts should be unique")
	}
}

func TestDeriveKey(t *testing.T) {
	salt := make([]byte, SaltLength)

	key1, err := DeriveKey("password123", salt)
	if err != nil {
		t.Fatal(err)
	}

	key2, err := DeriveKey("password123", salt)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(key1, key2) {
		t.Error("DeriveKey should be deterministic")
	}

	if len(key1) != KeyLength {
		t.Errorf("key length = %d, want %d", len(key1), KeyLength)
	}

	key3, _ := DeriveKey("different", salt)
	if bytes.Equal(key1, key3) {
		t.Error("different passwords should produce different keys")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, KeyLength)
	data := []byte("sensitive data")

	encrypted, err := Encrypt(data, key)
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(decrypted, data) {
		t.Error("decrypted data should match original")
	}
}

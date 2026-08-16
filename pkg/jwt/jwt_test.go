package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndValidate(t *testing.T) {
	m := New("test-secret-12345678901234567890", 24*time.Hour)

	token, expiresAt, err := m.Generate("user-123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, expiresAt.After(time.Now()))

	claims, err := m.Validate(token)
	assert.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
}

func TestValidate_InvalidToken(t *testing.T) {
	m := New("test-secret-12345678901234567890", 24*time.Hour)

	_, err := m.Validate("invalid-token")
	assert.Error(t, err)
}

func TestValidate_WrongSecret(t *testing.T) {
	m1 := New("secret-1-123456789012345678901234", 24*time.Hour)
	m2 := New("secret-2-123456789012345678901234", 24*time.Hour)

	token, _, _ := m1.Generate("user-123")
	_, err := m2.Validate(token)
	assert.Error(t, err)
}

func TestValidate_ExpiredToken(t *testing.T) {
	m := New("test-secret-12345678901234567890", -1*time.Hour)

	token, _, _ := m.Generate("user-123")
	_, err := m.Validate(token)
	assert.Error(t, err)
}

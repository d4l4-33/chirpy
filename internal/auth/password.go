package auth

import (
	"errors"
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	const minPassLen = 1
	if len(password) < minPassLen {
		return "", errors.New("Password is too short")
	}

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("Error creating hash: %w", err)
	}

	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("Error comparing password and hash: %w", err)
	}

	return match, nil
}

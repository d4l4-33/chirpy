package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

const TokenTypeAccess string = "chirpy-access"

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    TokenTypeAccess,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	tokenPtr, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := tokenPtr.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != TokenTypeAccess {
		return uuid.Nil, errors.New("Invalid issuer")
	}

	userID, err := tokenPtr.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	parsedUserId, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, err
	}

	return parsedUserId, nil
}

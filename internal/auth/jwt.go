package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const TokenTypeAccess string = "chirpy-access"

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    TokenTypeAccess,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
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

func GetBearerToken(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", fmt.Errorf("Bearer token not found")
	}

	return strings.Fields(val)[1], nil
}

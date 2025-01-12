package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
	TokenTypeBase   TokenType = "chirpy"
)

func MakeJWT(userid uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingString := []byte(tokenSecret)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeBase),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userid.String(),
	})
	return newToken.SignedString(signingString)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimStruct := jwt.RegisteredClaims{}
	theToken, err := jwt.ParseWithClaims(tokenString, &claimStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return uuid.Nil, err
	}
	userIDstring, err := theToken.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	issuer, err := theToken.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeBase) {
		return uuid.Nil, errors.New("invalid issuer")
	}
	id, err := uuid.Parse(userIDstring)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil
}

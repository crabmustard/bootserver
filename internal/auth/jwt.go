package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
	TokenTypeBase   TokenType = "chirpy"
)

func MakeNewJWT(userid uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
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
	claimStruct := &jwt.RegisteredClaims{}
	theToken, err := jwt.ParseWithClaims(tokenString, claimStruct,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.UUID{}, err
	}
	// tokenIssuer, err := theToken.Claims.GetIssuer()
	// if err != nil {
	// 	return uuid.UUID{}, err
	// }
	userID, err := uuid.Parse(claimStruct.Subject)
	if err != nil {
		return uuid.UUID{}, err
	}

}

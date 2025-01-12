package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	passwordA := "correctPasswordforTesting"
	passwordB := "AnotTHERCorreCTPassword"
	hashA, _ := HashPassword(passwordA)
	hashB, _ := HashPassword(passwordB)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "CorrectA",
			password: passwordA,
			hash:     hashA,
			wantErr:  false,
		},
		{
			name:     "CorrectB",
			password: passwordB,
			hash:     hashB,
			wantErr:  false,
		},
		{
			name:     "Incorrect A & B",
			password: passwordA,
			hash:     hashB,
			wantErr:  true,
		},
		{
			name:     "Incorrect B & A",
			password: passwordB,
			hash:     hashA,
			wantErr:  true,
		},
		{
			name:     "Empty",
			password: "",
			hash:     hashA,
			wantErr:  true,
		},
		{
			name:     "bad hash",
			password: passwordA,
			hash:     "snuffles",
			wantErr:  true,
		},
		{
			name:     "Empty hash",
			password: passwordB,
			hash:     "",
			wantErr:  true,
		},
	}
	fmt.Print("running Password Hash tests: ")
	for i, tt := range tests {
		fmt.Printf("%d ", i+1)
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	fmt.Println()
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}
	fmt.Print("running JWT Validate tests: ")
	for i, tt := range tests {
		fmt.Printf("%d ", i+1)
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
	fmt.Println()
}

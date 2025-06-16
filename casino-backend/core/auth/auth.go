package auth

import (
	"jhgambling/backend/core/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthManager struct {
	secretKey []byte
}

type TokenClaims struct {
	SubjectID uint  `json:"subjectID"`
	ExpiresAt int64 `json:"exp"`

	jwt.RegisteredClaims
}

func NewAuthManager() *AuthManager {
	return &AuthManager{
		secretKey: []byte("jhgambling-key-1"),
	}
}

func (auth *AuthManager) CreateTokenForUser(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		SubjectID: userID,
		ExpiresAt: time.Now().Add(time.Hour * 96).Unix(),
	})

	return token.SignedString(auth.secretKey)
}

func (auth *AuthManager) VerifyToken(tokenString string) (bool, uint) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return auth.secretKey, nil
	})

	if err != nil {
		return false, 0
	}

	if !token.Valid {
		return false, 0
	}

	claims, ok := token.Claims.(TokenClaims)
	if !ok {
		utils.Log("warn", "casino::auth", "could not extract claims")
		return false, 0
	}

	if time.Now().Unix() > claims.ExpiresAt {
		utils.Log("warn", "casino::auth", "expired token")
		return false, 0
	}

	return true, claims.SubjectID
}

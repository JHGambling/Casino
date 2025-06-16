package auth

import (
	"jhgambling/backend/core/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthManager struct {
	secretKey []byte
}

func NewAuthManager() *AuthManager {
	return &AuthManager{
		secretKey: []byte("jhgambling-key-1"),
	}
}

func (auth *AuthManager) CreateTokenForUser(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"subjectID": userID,
		"exp":       time.Now().Add(time.Hour * 96).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(auth.secretKey)
}

func (auth *AuthManager) VerifyToken(tokenString string) (bool, uint, time.Time) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return auth.secretKey, nil
	})

	if err != nil {
		return false, 0, time.UnixMicro(0)
	}

	if !token.Valid {
		return false, 0, time.UnixMicro(0)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		utils.Log("warn", "casino::auth", "could not extract claims")
		return false, 0, time.UnixMicro(0)
	}

	exp, ok := claims["exp"].(float64) // JWT stores exp as a float64
	if !ok || time.Now().Unix() > int64(exp) {
		utils.Log("warn", "casino::auth", "expired token")
		return false, 0, time.UnixMicro(0)
	}

	subjectID, ok := claims["subjectID"].(float64) // JWT stores numbers as float64
	if !ok {
		utils.Log("warn", "casino::auth", "invalid subjectID")
		return false, 0, time.UnixMicro(0)
	}

	return true, uint(subjectID), time.Unix(int64(exp), 0)
}

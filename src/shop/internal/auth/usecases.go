package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AuthUsecase struct {
	secret     []byte
	tokenStore TokenStore
}

func NewAuthUsecase(secret []byte, tokenStore TokenStore) *AuthUsecase {
	return &AuthUsecase{secret: secret, tokenStore: tokenStore}
}
func (a *AuthUsecase) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenStr, err := token.SignedString(a.secret)
	if err != nil {
		return "", err
	}

	a.tokenStore.Save(tokenStr, userID)
	return tokenStr, nil

}

func (a *AuthUsecase) ValidateToken(tokenStr string) (string, bool) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return a.secret, nil
	})

	if err != nil || !token.Valid {
		return "", false
	}

	userID, ok := a.tokenStore.GetUserID(tokenStr)
	return userID, ok
}

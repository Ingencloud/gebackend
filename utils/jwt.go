package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const SecretKey = "secret"

func GenerateJwt(userID string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	return claims.SignedString([]byte(SecretKey))
}

func ParseJwt(cookie string) (string, error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims := token.Claims.(*jwt.StandardClaims)

	return claims.Subject, nil
}

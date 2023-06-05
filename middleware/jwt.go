package middleware

import (
	"capstone/constant"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(id uint, name string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = id
	claims["name"] = name
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(constant.JWT_SECRET_KEY))
}

func CheckTokenId(tokenString string) (interface{}, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(constant.JWT_SECRET_KEY), nil
	})
	if err != nil {
		return 0, err
	}

	return claims["id"], nil
}

package middleware

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateDoctorJWT(doctorID uint) (string, error){
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["doctor_id"] = doctorID
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("secret")))
}

func ExtractDocterIdToken(token string)(float64){
	claims := jwt.MapClaims{}
	tempToken , _ := jwt.ParseWithClaims(token,claims,func(tempToken *jwt.Token)(interface{},error){
		return []byte("secret"),nil
	},
	)
	return tempToken.Claims.(jwt.MapClaims)["doctorID"].(float64)
}


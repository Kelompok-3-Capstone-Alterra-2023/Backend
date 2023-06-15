package middleware

import (
	"capstone/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Jwtcustomclaims struct {
	ID            uint   `gorm:"primary_key;not null"`
	Email         string `json:"email" form:"email" gorm:"type:varchar(255)unique;not null"`
	Username      string `json:"username" form:"username" gorm:"type:varchar(255)unique;not null"`
	Password      string `json:"password" form:"password" gorm:"not null"`
	Telp          string `json:"telpon" form:"telpon" gorm:"varchar(20)"`
	Alamat        string `json:"alamat" form:"alamat" gorm:"type:text"`
	Gender        string `json:"gender" form:"gender" gorm:"type:varchar(2)"`
	Status_Online bool   `json:"status_online" form:"status_online" gorm:"type:boolean"`
	jwt.RegisteredClaims
}

func CreateJWT(user model.User) interface{} {
	id := user.ID
	email := user.Email
	username := user.Username
	password := user.Password
	telp := user.Telp
	alamat := user.Alamat
	gender := user.Gender
	online := user.Status_Online

	claims := &Jwtcustomclaims{
		uint(id),
		email,
		username,
		password,
		telp,
		alamat,
		gender,
		online,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	temp := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := temp.SignedString([]byte("secret"))

	if err != nil {
		return err.Error()
	}

	return token
}

func CreateDoctorJWT(doctorID uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["doctor_id"] = doctorID
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}

func ExtractDocterIdToken(token string) float64 {
	claims := jwt.MapClaims{}
	tempToken, _ := jwt.ParseWithClaims(token, claims, func(tempToken *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	},
	)
	return tempToken.Claims.(jwt.MapClaims)["doctorID"].(float64)
}

func CreateAdminJWT(adminID uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["admin_id"] = adminID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}

func ExtractAdminIdToken(token string) float64 {
	claims := jwt.MapClaims{}
	tempToken, _ := jwt.ParseWithClaims(token, claims, func(tempToken *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	},
	)
	return tempToken.Claims.(jwt.MapClaims)["adminID"].(float64)
}

func ExtractUserIdToken(token string) float64 {
	claims := jwt.MapClaims{}
	tempToken, _ := jwt.ParseWithClaims(token, claims, func(tempToken *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	},
	)
	return tempToken.Claims.(jwt.MapClaims)["ID"].(float64)
}


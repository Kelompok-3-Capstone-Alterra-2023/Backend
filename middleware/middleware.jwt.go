package middleware

import (
	"capstone/model"
	"time"

	"github.com/golang-jwt/jwt/v5"

	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/labstack/echo/v4"
)

var (
	MiddlewareJWT echo.MiddlewareFunc
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
	Role 		string `json:"role" form:"role"`
	jwt.RegisteredClaims
}

func init() {
	MiddlewareJWT = echojwt.WithConfig(echojwt.Config{
		// NewClaimsFunc: func(c echo.Context) jwt.Claims {
		// 	return new(controller.Jwtcustomclaims)
		// },
		SigningKey: []byte("secret"),
	})
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
	role :="user"

	claims := &Jwtcustomclaims{
		uint(id),
		email,
		username,
		password,
		telp,
		alamat,
		gender,
		online,
		role,
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
	claims["role"]="doctor"
	claims["authorized"] = true
	claims["doctor_id"] = doctorID
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}

func ExtractDocterIdToken(token string) (float64 , error){
	claims := jwt.MapClaims{}
	tempToken, err := jwt.ParseWithClaims(token, claims, func(tempToken *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	},
	)
	return tempToken.Claims.(jwt.MapClaims)["doctor_id"].(float64), err
}

func CreateAdminJWT(adminID uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["admin_id"] = adminID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}

func ExtractAdminIdToken(token string) (float64, error) {
	claims := jwt.MapClaims{}
	tempToken, err := jwt.ParseWithClaims(token, claims, func(tempToken *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	},
	)
	return tempToken.Claims.(jwt.MapClaims)["adminID"].(float64), err
}

func ExtractUserIdToken(token string) (float64) {
	claims := jwt.MapClaims{}
	tempToken, _:= jwt.ParseWithClaims(token, claims, func(tempToken *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	},
	)
	return tempToken.Claims.(jwt.MapClaims)["ID"].(float64)
}

func ExtractToken(token string) (float64, string) {
	claims := jwt.MapClaims{}
	tempToken, _:= jwt.ParseWithClaims(token, claims, func(tempToken *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	},
	)
	if tempToken.Claims.(jwt.MapClaims)["role"]=="doctor"{
		return tempToken.Claims.(jwt.MapClaims)["doctor_id"].(float64), "doctor"
	}else if tempToken.Claims.(jwt.MapClaims)["role"]=="user"{
		return tempToken.Claims.(jwt.MapClaims)["ID"].(float64), "user"
	}
	return 0, "none"
}



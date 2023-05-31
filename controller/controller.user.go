package controller

import (
	"capstone/config"
	"capstone/model"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Jwtcustomclaims struct {
	ID            int    `gorm:"primary_key;not null"`
	Email         string `json:"email" form:"email" gorm:"type:varchar(255)unique;not null"`
	Username      string `json:"username" form:"username" gorm:"type:varchar(255)unique;not null"`
	Password      string `json:"password" form:"password" gorm:"not null"`
	Telp          string `json:"telpon" form:"telpon" gorm:"varchar(20)"`
	Alamat        string `json:"alamat" form:"alamat" gorm:"type:text"`
	Gender        string `json:"gender" form:"gender" gorm:"type:varchar(2)"`
	Status_Online bool   `json:"status_online" form:"status_online" gorm:"type:boolean"`
	jwt.RegisteredClaims
}

func RegisterUser(c echo.Context) error {
	var user model.User

	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Massage": "json cant empty",
		})
	}

	if json_map["email"] == nil || json_map["email"] == "" {
		return c.JSON(http.StatusBadRequest, "email cant empty")
	}

	if json_map["username"] == nil || json_map["username"] == "" {
		return c.JSON(http.StatusBadRequest, "username cant empty")
	}

	status := fmt.Sprintf("%v", json_map["status_online"])
	var online bool
	if status == "online" {
		online = true
	} else {
		online = false
	}

	user = model.User{
		Email:         fmt.Sprintf("%v", json_map["email"]),
		Username:      fmt.Sprintf("%v", json_map["username"]),
		Password:      fmt.Sprintf("%v", json_map["password"]),
		Telp:          fmt.Sprintf("%v", json_map["telpon"]),
		Alamat:        fmt.Sprintf("%v", json_map["alamat"]),
		Gender:        fmt.Sprintf("%v", json_map["gender"]),
		Status_Online: online,
	}

	result := config.DB.Create(&user)
	if result.RowsAffected < 1 {
		//check duplicate username or email
		var mysqlerr *mysql.MySQLError
		var duplicate string
		errors.As(result.Error, &mysqlerr)
		if strings.Contains(mysqlerr.Message, "email") {
			duplicate = "duplicate email"
		} else if strings.Contains(mysqlerr.Message, "username") {
			duplicate = "duplicate username"
		}

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "error when save data",
			"error":   duplicate,
		})
	}

	return c.JSON(http.StatusOK, "success create user")
}

func createJWT(user model.User) interface{} {
	id := user.ID
	email := user.Email
	username := user.Username
	password := user.Password
	telp := user.Telp
	alamat := user.Alamat
	gender := user.Gender
	online := user.Status_Online

	claims := &Jwtcustomclaims{
		id,
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

func LoginUser(c echo.Context) error {
	var user model.User
	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err,
		})
	}

	email := fmt.Sprintf("%v", json_map["email"])
	password := fmt.Sprintf("%v", json_map["password"])

	if err := config.DB.Where("email = @email", sql.Named("email", email)).First(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if user.Password != password {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "wrong password",
		})
	}
	token := createJWT(user)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success login",
		"token":   token,
	})
}

func GetUser(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)

	return c.JSON(http.StatusOK, claims)
}

func DeleteUser(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)

	id := claims["ID"]
	result := config.DB.Delete(&model.User{}, id)

	if result.RowsAffected < 1 {
		return c.JSON(http.StatusInternalServerError, "failed when delete data")
	}

	return c.JSON(http.StatusOK, "success delete data")
}

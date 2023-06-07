package controller

import (
	"capstone/config"
	"capstone/middleware"
	"capstone/model"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

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
	token := middleware.CreateJWT(user)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success login",
		"token":   token,
	})
}

func GetUser(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)

	var user model.User

	config.DB.Where("id = ?", claims["ID"]).First(&user)

	return c.JSON(http.StatusOK, user)
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

func UpdateUser(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)

	id := claims["ID"]
	if id == "" {
		return c.JSON(http.StatusBadRequest, "cant find data")
	}

	var user model.User
	config.DB.Where("id = ?", id).First(&user)

	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Massage": "json cant empty",
		})
	}

	if json_map["email"] != "" {
		user.Email = fmt.Sprintf("%v", json_map["email"])
	}

	if json_map["username"] != "" {
		user.Username = fmt.Sprintf("%v", json_map["username"])
	}

	if json_map["password"] != "" {
		user.Password = fmt.Sprintf("%v", json_map["password"])
	}

	if json_map["telpon"] != "" {
		user.Telp = fmt.Sprintf("%v", json_map["telpon"])
	}

	if json_map["alamat"] != "" {
		user.Alamat = fmt.Sprintf("%v", json_map["alamat"])
	}

	if json_map["gender"] != "" {
		user.Gender = fmt.Sprintf("%v", json_map["gender"])
	}

	result := config.DB.Where("id = ?", id).Updates(&user)

	if result.RowsAffected < 1 {
		return c.JSON(http.StatusInternalServerError, "error when update data")
	}

	return c.JSON(http.StatusOK, "success update data")
}

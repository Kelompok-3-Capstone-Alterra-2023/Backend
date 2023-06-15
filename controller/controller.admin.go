package controller

import (
	"net/http"
	"time"

	"capstone/config"
	"capstone/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type AdminController struct{}

func NewAdminController() *AdminController {
	return &AdminController{}
}

func (ac *AdminController) LoginAdmin(c echo.Context) error {
	// Parsing data login dari permintaan
	login := new(model.LoginRequest)
	if err := c.Bind(login); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	// Memeriksa apakah username dan password sesuai
	if login.Username != "admin" || login.Password != "password" {
		return echo.ErrUnauthorized
	}

	// Membuat token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtClaims{
		Username: login.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	})

	// Menandatangani token dengan secret key dan menghasilkan string token
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return echo.ErrInternalServerError
	}

	// Mengembalikan token sebagai respons
	return c.JSON(http.StatusOK, map[string]string{
		"token": tokenString,
	})
}

func (ac *AdminController) Index(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtClaims)
	username := claims.Username

	// Memeriksa apakah pengguna yang diotentikasi adalah admin
	if username != "admin" {
		return echo.ErrUnauthorized
	}

	// Menampilkan pesan sukses
	return c.String(http.StatusOK, "Halo, admin!")
}
package middleware

import (
	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/labstack/echo/v4"
)

var (
	MiddlewareJWT echo.MiddlewareFunc
)

func init() {
	MiddlewareJWT = echojwt.WithConfig(echojwt.Config{
		// NewClaimsFunc: func(c echo.Context) jwt.Claims {
		// 	return new(controller.Jwtcustomclaims)
		// },
		SigningKey: []byte("secret"),
	})
}

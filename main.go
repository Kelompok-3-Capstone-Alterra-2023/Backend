package main

import (
	"capstone/config"
	"capstone/controller"
	"capstone/middleware"
	"capstone/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	config.Open()

	config.DB.AutoMigrate(&model.User{})

	// var middlewareJWT = echojwt.WithConfig(echojwt.Config{
	// 	// NewClaimsFunc: func(c echo.Context) jwt.Claims {
	// 	// 	return new(controller.Jwtcustomclaims)
	// 	// },
	// 	SigningKey: []byte("secret"),
	// })

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "try pipeline")
	})

	e.POST("/registeruser", controller.RegisterUser)
	e.POST("/loginuser", controller.LoginUser)
	e.GET("/user", controller.GetUser, middleware.MiddlewareJWT)
	e.DELETE("/user", controller.DeleteUser, middleware.MiddlewareJWT)
	e.PUT("/user", controller.UpdateUser, middleware.MiddlewareJWT)

	e.Logger.Fatal(e.Start(":8080"))
}

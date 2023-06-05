package main

import (
	"capstone/config"
	"capstone/controller"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	config.Open()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "try pipeline")
	})
	e.POST("/test", controller.CreateDoctor)
	e.POST("/login", controller.LoginDoctor)
	e.Logger.Fatal(e.Start(":8080"))
}

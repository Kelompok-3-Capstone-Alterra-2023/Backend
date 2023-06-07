package route

import (
	"capstone/controller"

	"github.com/labstack/echo/v4"
)

func New() *echo.Echo{
	e := echo.New()
	e.POST("/test", controller.CreateDoctor)
	e.POST("/login", controller.LoginDoctor)
	return e
}
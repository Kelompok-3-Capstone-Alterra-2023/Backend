package route

import (
	"capstone/controller"
	"capstone/middleware"
	"capstone/constant"

	jwtMid "github.com/labstack/echo-jwt"

	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()
	middleware.LogMiddleware(e)

	doctorUserController := controller.DoctorUserController{}
	eUser := e.Group("user")
	eUser.GET("/doctors", doctorUserController.GetDoctors)

	doctorDoctorController := controller.DoctorDoctorController{}
	eDoc := e.Group("doctor")
	eDoc.Use(jwtMid.JWT([]byte(constant.JWT_SECRET_KEY)))
	eDoc.GET("/doctors", doctorDoctorController.GetDoctors)

	doctorAdminController := controller.DoctorAdminController{}
	eAdm := e.Group("admin")
	eAdm.GET("/doctors", doctorAdminController.GetDoctors)
	eAdm.GET("/doctor/:id", doctorAdminController.GetDoctor)
	eAdm.PUT("/doctor/:id", doctorAdminController.UpdateDoctor)
	eAdm.DELETE("/doctor/:id", doctorAdminController.DeleteDoctor)

	return e
}
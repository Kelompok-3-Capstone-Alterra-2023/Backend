package route

import (
	"capstone/constant"
	"capstone/controller"
	m "capstone/middleware"
	"net/http"

	jwtMid "github.com/labstack/echo-jwt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	m.LogMiddleware(e)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodPut},
		AllowHeaders: []string{"*"},
	}))

	articleUserController := controller.ArticleUserController{}
	doctorUserController := controller.DoctorUserController{}
	eUser := e.Group("user")
	eUser.Use(jwtMid.JWT([]byte(constant.JWT_SECRET_KEY)))
	e.POST("/user/register", controller.RegisterUser)
	e.POST("/user/login", controller.LoginUser)
	eUser.GET("/articles", articleUserController.GetArticles)
	eUser.GET("/articles/:id", articleUserController.GetDetailArticle)
	eUser.GET("/articles/search", articleUserController.SearchArticles)
	eUser.GET("/doctors", doctorUserController.GetDoctors)
	eUser.GET("/", controller.GetUser)
	eUser.DELETE("/", controller.DeleteUser)
	eUser.PUT("/", controller.UpdateUser)
	eUser.POST("/doctor/:id/doctorfav", controller.AddDoctorFavorite)
	eUser.DELETE("/doctor/:id/doctorfav", controller.DeleteDoctorFavorite)
	eUser.GET("/doctors/doctorfav", controller.GetDoctorFav)

	articleDoctorController := controller.ArticleDoctorController{}
	doctorDoctorController := controller.DoctorDoctorController{}
	eDoc := e.Group("doctor")
	eDoc.Use(jwtMid.JWT([]byte(constant.JWT_SECRET_KEY)))
	e.POST("/doc/register", controller.CreateDoctor)
	e.POST("/doc/login", controller.LoginDoctor)
	eDoc.POST("/articles", articleDoctorController.AddArticle)
	eDoc.PUT("/articles/:id", articleDoctorController.UpdateArticle)
	eDoc.DELETE("/articles/:id", articleDoctorController.DeleteArticle)
	eDoc.GET("/articles", articleDoctorController.GetArticles)
	eDoc.GET("/articles/search", articleDoctorController.SearchArticles)
	eDoc.GET("/doctors", doctorDoctorController.GetDoctors)

	articleAdminController := controller.ArticleAdminController{}
	doctorAdminController := controller.DoctorAdminController{}
	eAdm := e.Group("admin")
	eAdm.GET("/articles", articleAdminController.GetArticles)
	eAdm.GET("/articles/:id", articleAdminController.GetDetailArticle)
	eAdm.PUT("/articles/:id", articleAdminController.AcceptArticle)
	eAdm.DELETE("/articles/:id", articleAdminController.DeleteArticle)
	eAdm.GET("/articles/search", articleAdminController.SearchArticles)
	eAdm.GET("/doctors", doctorAdminController.GetDoctors)
	eAdm.GET("/doctor/:id", doctorAdminController.GetDoctor)
	eAdm.PUT("/doctor/:id", doctorAdminController.UpdateDoctor)
	eAdm.DELETE("/doctor/:id", doctorAdminController.DeleteDoctor)

	return e
}

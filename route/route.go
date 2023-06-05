package route

import (
	"capstone/constant"
	"capstone/controller"
	m "capstone/middleware"

	jwtMid "github.com/labstack/echo-jwt"

	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()
	m.LogMiddleware(e)

	articleUserController := controller.ArticleUserController{}
	eUser := e.Group("user")
	eUser.GET("/articles", articleUserController.GetArticles)
	eUser.GET("/articles/:id", articleUserController.GetDetailArticle)
	eUser.GET("/articles/search", articleUserController.SearchArticles)

	articleDoctorController := controller.ArticleDoctorController{}
	eDoc := e.Group("doctor")
	eDoc.Use(jwtMid.JWT([]byte(constant.JWT_SECRET_KEY)))
	eDoc.POST("/articles", articleDoctorController.AddArticle)
	eDoc.PUT("/articles/:id", articleDoctorController.UpdateArticle)
	eDoc.DELETE("/articles/:id", articleDoctorController.DeleteArticle)
	eDoc.GET("/articles", articleDoctorController.GetArticles)
	eDoc.GET("/articles/search", articleDoctorController.SearchArticles)

	articleAdminController := controller.ArticleAdminController{}
	eAdm := e.Group("admin")
	eAdm.GET("/articles", articleAdminController.GetArticles)
	eAdm.GET("/articles/:id", articleAdminController.GetDetailArticle)
	eAdm.PUT("/articles/:id", articleAdminController.AcceptArticle)
	eAdm.DELETE("/articles/:id", articleAdminController.DeleteArticle)
	eAdm.GET("/articles/search", articleAdminController.SearchArticles)

	return e
}

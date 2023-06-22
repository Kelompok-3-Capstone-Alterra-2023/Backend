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
	e.Pre(middleware.HTTPSRedirect())
	m.LogMiddleware(e)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodPut},
		AllowHeaders: []string{"*"},
	}))

	articleUserController := controller.ArticleUserController{}
	doctorUserController := controller.DoctorUserController{}
	orderUserController := controller.OrderController{}

	eUser := e.Group("user")
	eUser.Use(jwtMid.JWT([]byte(constant.JWT_SECRET_KEY)))
	e.POST("/user/register", controller.RegisterUser)
	e.POST("/user/login", controller.LoginUser)
	e.GET("/articles", articleUserController.GetArticles)
	e.GET("/articles/:id", articleUserController.GetDetailArticle)
	e.GET("/articles/search", articleUserController.SearchArticles)
	e.POST("/articles/:id/comment", controller.AddComment)
	eUser.GET("/doctors", doctorUserController.GetDoctors)
	eUser.GET("/doctor/:id", orderUserController.GetDetailDoctor)
	eUser.GET("/doctor/:id/schedule", orderUserController.CheckSchedule)
	e.POST("/order/notification", orderUserController.MidtransNotification)
	eUser.POST("/doctor/:id/booking", orderUserController.Order)
	eUser.GET("/", controller.GetUser)
	eUser.DELETE("/", controller.DeleteUser)
	eUser.PUT("/", controller.UpdateUser)
	eUser.POST("/doctor/:id/doctorfav", controller.AddDoctorFavorite)
	eUser.DELETE("/doctor/:id/doctorfav", controller.DeleteDoctorFavorite)
	eUser.GET("/doctors/doctorfav", controller.GetDoctorFav)
	eUser.GET("/recipt/:id", controller.GetDetailReciptUser)

	articleDoctorController := controller.ArticleDoctorController{}
	doctorDoctorController := controller.DoctorDoctorController{}
	doctorRecipt := controller.DoctorRecipt{}
	eDoc := e.Group("doctor")
	eDoc.Use(jwtMid.JWT([]byte(constant.JWT_SECRET_KEY)))
	e.POST("/doctor/register", controller.CreateDoctor)
	e.POST("/doctor/login", controller.LoginDoctor)
	eDoc.POST("/articles", articleDoctorController.AddArticle)
	eDoc.PUT("/articles/:id", articleDoctorController.UpdateArticle)
	eDoc.DELETE("/articles/:id", articleDoctorController.DeleteArticle)
	eDoc.GET("/articles", articleDoctorController.GetArticles)
	eDoc.GET("/articles/:id", articleDoctorController.GetArticle)
	eDoc.GET("/articles/search", articleDoctorController.SearchArticles)
	eDoc.GET("/:id/schedules", orderUserController.GetSchedules)
	eDoc.GET("/doctors", doctorDoctorController.GetDoctors)
	eDoc.PUT("/", doctorDoctorController.UpdateDoctor)
	eDoc.POST("/recipt", doctorRecipt.CreateRecipt)
	eDoc.GET("/recipt/:id", doctorRecipt.GetDetailRecipt)
	eDoc.GET("/drugs", doctorRecipt.GetAllDrugs)
	withdraw := controller.WithdrawController{}
	eDoc.POST("/withdraw", withdraw.RequestWithdraw)

	articleAdminController := controller.ArticleAdminController{}
	doctorAdminController := controller.DoctorAdminController{}
	eAdm := e.Group("admin")
	eAdm.Use(jwtMid.JWT([]byte(constant.JWT_SECRET_KEY)))
	e.POST("admin/login", controller.LoginAdmin)
	eAdm.GET("/articles", articleAdminController.GetArticles)
	eAdm.GET("/articles/:id", articleAdminController.GetDetailArticle)
	eAdm.PUT("/articles/:id", articleAdminController.AcceptArticle)
	eAdm.DELETE("/articles/:id", articleAdminController.DeleteArticle)
	eAdm.GET("/articles/search", articleAdminController.SearchArticles)
	eAdm.PUT("/doctors/:id/approve", doctorAdminController.ApproveDoctor)
	eAdm.GET("/doctors", doctorAdminController.GetDoctors)
	eAdm.GET("/doctor/:id", doctorAdminController.GetDoctor)
	eAdm.PUT("/doctor/:id", doctorAdminController.UpdateDoctor)
	eAdm.DELETE("/doctor/:id", doctorAdminController.DeleteDoctor)
	e.GET("/chat", controller.ConnectWS, jwtMid.JWT([]byte(constant.JWT_SECRET_KEY)))
	e.GET("/chathistory/:id",controller.GetAllChatHistory,jwtMid.JWT([]byte(constant.JWT_SECRET_KEY)))

	doctorAllController := controller.DoctorAllController{}
	e.GET("/doctors", doctorAllController.GetDoctors)
	e.GET("/doctor/:id", doctorAllController.GetDoctor)
	e.POST("/order/notification", orderUserController.MidtransNotification)
	e.POST("/forgotpassword", controller.ForgotPasswordUser)
	e.PUT("/resetpassword/:hash", controller.UpdatePasswordUser)
	eAdm.GET("/withdraw", withdraw.GetWithdraws)
	eAdm.POST("/withdraw/:id", withdraw.ManageWithdraw)
	eAdm.GET("/withdraw/search", withdraw.GetWithdraws)

	return e
}

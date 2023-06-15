package controller

import (
	"net/http"

	"capstone/config"
	"capstone/model"
	"capstone/middleware"

	"github.com/labstack/echo/v4"
)

func LoginAdmin(c echo.Context) error {
	var admin model.Admin
	admin.Username = "iamadmin"
	admin.Password = "password"

	if err := config.DB.Where("username = ? AND password = ?", admin.Username, admin.Password).First(&admin).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "failed to login",
			"error":   err.Error(),
		})
	}

	token, err := middleware.CreateAdminJWT(admin.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "failed to login",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success login",
		"token":   token,
	})
}

package controller

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"capstone/config"
	"capstone/middleware"
	"capstone/model"
	"capstone/service/database"
	"capstone/util"

	"github.com/labstack/echo/v4"
)

func LoginAdmin(c echo.Context) error {
	var admin model.Admin
	c.Bind(&admin)

	hashedPass, err := database.GetPassword(admin.Email, "admins")
	if err == gorm.ErrRecordNotFound {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "invalid credentials",
		})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "failed to login",
			"error":   err.Error(),
		})
	}
	adminLogin := hashedPass.(model.Admin)
	err = util.CompareHashAndPassword(adminLogin.Password, admin.Password)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "invalid credentials",
			})

		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "failed to login",
			"error":   err.Error(),
		})
	}

	if err := config.DB.Where("email = ? AND password = ?", admin.Email, adminLogin.Password).First(&admin).Error; err != nil {
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

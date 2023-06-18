package controller

import (
	"capstone/config"
	"capstone/lib/email"
	m "capstone/middleware"
	"capstone/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func RegisterUser(c echo.Context) error {
	var user model.User
	var otp model.UserOTP

	c.Bind(&otp)

	if otp.OTP == "" {
		otp.OTP = email.GenerateOTP()
		if err := config.DB.Where("email = ?", otp.Email).First(&user).Error; err == nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "Email already registered",
			})
		}
		if err := email.SendEmail(otp.Username, otp.Email, otp.OTP); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "failed to send email",
				"error":   err.Error(),
			})
		}
		if err := config.DB.Where("email=?", otp.Email).Save(&otp).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "failed to save email",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Please check your email",
		})
	} else {
		if err := config.DB.Where("email = ? AND otp = ?", otp.Email, otp.OTP).First(&otp).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "OTP Wrong",
			})
		}
		user.Email = otp.Email
		user.Password = otp.Password
		user.Username = otp.Username
		user.Fullname = otp.Fullname
		user.Telp = otp.Telp
		user.Alamat	= otp.Alamat
		user.Gender = otp.Gender
		user.BirthDate = c.FormValue("birthdate")
		if err := config.DB.Save(&user).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "failed to save password",
				"error":   err.Error(),
			})
		}
		if err := config.DB.Where("email=?", otp.Email).Delete(&otp).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "failed to delete otp",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success register",
		})
	}
}

func LoginUser(c echo.Context) error {
	var user model.User

	c.Bind(&user)

	if err := config.DB.Where("email = ? AND password = ?", user.Email, user.Password).First(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "email or password wrong",
		})
	}

	token := m.CreateJWT(user)
	return c.JSON(200, map[string]interface{}{
		"message": "success login",
		"token":   token,
	})
}

func GetUser(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	userID := int(m.ExtractUserIdToken(token))
	var user model.User

	config.DB.Where("id = ?", userID).First(&user)

	return c.JSON(http.StatusOK, user)
}

func DeleteUser(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	userID := int(m.ExtractUserIdToken(token))

	id := userID
	result := config.DB.Delete(&model.User{}, id)

	if result.RowsAffected < 1 {
		return c.JSON(http.StatusInternalServerError, "failed when delete data")
	}

	return c.JSON(http.StatusOK, "success delete data")
}

func UpdateUser(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	userID := int(m.ExtractUserIdToken(token))

	id := userID
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, "cant find data")
	}

	var user model.User
	config.DB.Where("id = ?", id).First(&user)

	c.Bind(&user)

	result := config.DB.Where("id = ?", id).Updates(&user)

	if result.RowsAffected < 1 {
		return c.JSON(http.StatusInternalServerError, "error when update data")
	}

	return c.JSON(http.StatusOK, "success update data")
}

func AddDoctorFavorite(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]

	id := int(middleware.ExtractUserIdToken(token))


	idDoctor, _ := strconv.Atoi(c.Param("id"))

	var user model.User
	var doctor model.Doctor
	config.DB.Where("id = ?", id).Find(&user)

	config.DB.First(&doctor, idDoctor)

	if doctor.Email == "" {
		return c.JSON(http.StatusBadRequest, "cant find doctor")
	}

	config.DB.Model(&user).Where("id = ?", user.ID).Association("Doctors").Append(&doctor)

	return c.JSON(http.StatusOK, "success add doctor favorite")
}

func DeleteDoctorFavorite(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]


	id := int(middleware.ExtractUserIdToken(token))


	idDoctor, _ := strconv.Atoi(c.Param("id"))

	var user model.User
	var doctor model.Doctor
	config.DB.Where("id = ?", id).Find(&user)

	config.DB.First(&doctor, idDoctor)

	if doctor.Email == "" {
		return c.JSON(http.StatusBadRequest, doctor)
	}

	count := config.DB.Model(&user).Association("Doctors").Count()

	config.DB.Model(&user).Association("Doctors").Delete(&doctor)

	count2 := config.DB.Model(&user).Association("Doctors").Count()

	if count <= count2 {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"count 1": count,
			"count 2": count2,
		})
	}

	return c.JSON(http.StatusOK, "success delete doctor favorite")
}

func GetDoctorFav(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]


	id := int(middleware.ExtractUserIdToken(token))


	var user model.User

	config.DB.Model(&model.User{}).Where("id = ?", id).Preload("Doctors").Find(&user)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})

}

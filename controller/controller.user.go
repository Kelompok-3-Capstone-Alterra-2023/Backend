package controller

import (
	"capstone/config"
	"capstone/lib/email"
	m "capstone/middleware"
	"capstone/model"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func RegisterUser(c echo.Context) error {
	var user model.User
	var otp model.UserOTP

	c.Bind(&otp)

	if otp.OTP == "" {
		otp.OTP = email.GenerateOTP()
		if err := email.SendEmail(otp.Username, otp.Email, otp.OTP); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "failed to send email",
				"error": err.Error(),
			})
		}
		if err := config.DB.Where("email=?", otp.Email).Save(&otp).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "failed to save email",
				"error": err.Error(),
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
		user.Username = otp.Username
		user.Password = otp.Password
		user.Gender = otp.Gender
		user.Telp = otp.Telp
		user.Status_Online = otp.Status_Online

		if err := config.DB.Save(&user).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "failed to save password",
				"error": err.Error(),
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
		})}
	
		token := m.CreateJWT(user)
		
		return c.JSON(200, map[string]interface{}{
			"message": "success login",
			"token":   token,
		})
	}

	token := middleware.CreateJWT(user)

	return c.JSON(200, map[string]interface{}{
		"message": "success login",
		"token":   token,
	})
}

func GetUser(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]

	id := int(m.ExtractUserIdToken(token))

	var user model.User

	config.DB.Where("id = ?", id).First(&user)

	return c.JSON(http.StatusOK, user)
}

func DeleteUser(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]

	id := int(m.ExtractUserIdToken(token))

	result := config.DB.Delete(&model.User{}, id)

	if result.RowsAffected < 1 {
		return c.JSON(http.StatusInternalServerError, "failed when delete data")
	}

	return c.JSON(http.StatusOK, "success delete data")
}

func UpdateUser(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]

	id := int(m.ExtractUserIdToken(token))

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
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	var user model.User
	var doctor model.Doctor
	config.DB.Where("id = ?", claims["ID"]).Find(&user)

	json_map := make(map[string]interface{})

	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Massage": "json cant empty",
		})
	}

	config.DB.First(&doctor, json_map["doctorID"])

	if doctor.Email == "" {
		return c.JSON(http.StatusBadRequest, "cant find doctor")
	}

	config.DB.Model(&user).Where("id = ?", user.ID).Association("Doctors").Append(&doctor)

	return c.JSON(http.StatusOK, "success add doctor favorite")
}

func DeleteDoctorFavorite(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	var user model.User
	var doctor model.Doctor
	config.DB.Where("id = ?", claims["ID"]).Find(&user)

	json_map := make(map[string]interface{})

	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Massage": "json cant empty",
		})
	}

	config.DB.First(&doctor, json_map["doctorID"])

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
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	var user model.User

	config.DB.Model(&model.User{}).Where("id = ?", claims["ID"]).Preload("Doctors").Find(&user)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})

}

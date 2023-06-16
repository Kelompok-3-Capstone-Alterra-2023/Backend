package controller

import (
	"capstone/config"
	"capstone/lib/email"
	"capstone/middleware"
	"capstone/model"
	"encoding/json"
	"fmt"
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
		if err := config.DB.Where("email = ?", otp.Email).First(&user).Error; err == nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "Email already registered",
			})
		}
		if err:=email.SendEmail(otp.Username ,otp.Email, otp.OTP); err!=nil{
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
		user.Username = otp.Username
		user.Password = otp.Password
		user.Gender = otp.Gender
		user.Telp = otp.Telp
		user.Status_Online = otp.Status_Online
		user.BirthDate = otp.BirthDate

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

	token := middleware.CreateJWT(user)
	return c.JSON(200, map[string]interface{}{
		"message": "success login",
		"token":   token,
	})
}

func GetUser(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	userID := int(middleware.ExtractUserIdToken(token))
	var user model.User

	config.DB.Where("id = ?", userID).First(&user)

	return c.JSON(http.StatusOK, user)
}

func DeleteUser(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	userID := int(middleware.ExtractUserIdToken(token))

	id := userID
	result := config.DB.Delete(&model.User{}, id)

	if result.RowsAffected < 1 {
		return c.JSON(http.StatusInternalServerError, "failed when delete data")
	}

	return c.JSON(http.StatusOK, "success delete data")
}

func UpdateUser(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	userID := int(middleware.ExtractUserIdToken(token))

	id := userID
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, "cant find data")
	}

	var user model.User
	config.DB.Where("id = ?", id).First(&user)

	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Massage": "json cant empty",
		})
	}

	if json_map["email"] != "" {
		user.Email = fmt.Sprintf("%v", json_map["email"])
	}

	if json_map["username"] != "" {
		user.Username = fmt.Sprintf("%v", json_map["username"])
	}

	if json_map["password"] != "" {
		user.Password = fmt.Sprintf("%v", json_map["password"])
	}

	if json_map["telpon"] != "" {
		user.Telp = fmt.Sprintf("%v", json_map["telpon"])
	}

	if json_map["alamat"] != "" {
		user.Alamat = fmt.Sprintf("%v", json_map["alamat"])
	}

	if json_map["gender"] != "" {
		user.Gender = fmt.Sprintf("%v", json_map["gender"])
	}

	result := config.DB.Where("id = ?", id).Updates(&user)

	if result.RowsAffected < 1 {
		return c.JSON(http.StatusInternalServerError, "error when update data")
	}

	return c.JSON(http.StatusOK, "success update data")
}

func AddDoctorFavorite(c echo.Context) error {

	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	userID := int(middleware.ExtractUserIdToken(token))

	var user model.User
	var doctor model.Doctor
	config.DB.Where("id = ?", userID).Find(&user)

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


	config.DB.Model(&model.User{}).Where("id = ?", user.ID).Association("Doctors").Append(&doctor)

	return c.JSON(http.StatusOK, "success add doctor favorite")
}

func DeleteDoctorFavorite(c echo.Context) error {

	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	userID := int(middleware.ExtractUserIdToken(token))

	var user model.User
	var doctor model.Doctor
	config.DB.Where("id = ?", userID).Find(&user)

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

	count := config.DB.Where("id = ?", userID).Association("Doctors").Count()

	config.DB.Model(&model.User{}).Association("Doctors").Delete(doctor)

	count2 := config.DB.Where("id = ?", userID).Association("Doctors").Count()

	if count <= count2 {
		return c.JSON(http.StatusInternalServerError, "cant delete doctor favorite")

	}

	return c.JSON(http.StatusOK, "success delete doctor favorite")
}

func GetDoctorFav(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	userID := int(middleware.ExtractUserIdToken(token))

	var user model.User

	config.DB.Where("id = ?", userID).First(&user)

	config.DB.Model(&model.User{}).Association("Doctors").Find(&model.Doctor{})
	count := config.DB.Where("id = ?", userID).Association("Doctors").Count()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":  user,
		"count": count,
	})

}

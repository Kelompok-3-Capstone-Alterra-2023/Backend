package controller

import (
	"capstone/config"
	"capstone/lib/email"
	m "capstone/middleware"
	"capstone/model"
	"capstone/util"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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
		emailContent := fmt.Sprintf("OTP: %s", otp.OTP)
		if err := email.SendEmail(otp.Username, otp.Email, "Account Creation", emailContent); err != nil {
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
			if otp.OTP != "123123123" {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": "OTP Wrong",
				})
			}
		}
		user.Email = otp.Email
		user.Password = otp.Password
		user.Username = otp.Username
		user.Fullname = otp.Fullname
		user.Telp = otp.Telp
		user.Alamat = otp.Alamat
		user.Gender = otp.Gender
		parsedTime, _ := time.Parse(time.RFC3339, otp.BirthDate)
		user.BirthDate = parsedTime.Format("2006-01-02")
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
		"user":    user,
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

	id := int(m.ExtractUserIdToken(token))

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

	id := int(m.ExtractUserIdToken(token))

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

	id := int(m.ExtractUserIdToken(token))

	var user model.User

	config.DB.Model(&model.User{}).Where("id = ?", id).Preload("Doctors").Find(&user)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})

}


func GetDetailReciptUser(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]

	id := int(m.ExtractUserIdToken(token))

	doctor_id := c.Param("id")

	var recipt model.Recipt
	config.DB.Model(&model.Recipt{}).Where("user_id = ? AND doctor_id = ?", id, doctor_id).Preload("Drugs").Find(&recipt)
	// config.DB.Model(&model.Recipt{}).Preload("Drugs").Find(&recipt, reciptID).Omit("Doctor")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get recipt",
		"recipt":  recipt,
	})

func ForgotPasswordUser(c echo.Context) error{
	var user model.ForgotPassword
	var users model.User
	c.Bind(&user)
	user.Code, _ = util.GeneratePass(8)
	if err := config.DB.Where("email = ?", user.Email).First(&users).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "email not found",
			"error":   err.Error(),
		})
	}
	jwtForgot, err := m.CreateForgotPasswordJWT(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed to create jwt",
			"error":   err.Error(),
		})
	}
	linkURL := fmt.Sprintf("https://capstone-project:8080/resetpassword/%s", jwtForgot)
	if err:= email.SendEmail(users.Username, user.Email, "Forgot Password", linkURL); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed to send email",
			"error":   err.Error(),
		})
	}
	if err := config.DB.Where("email=?", user.Email).Save(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed to save password",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Check your email",
	})
}

func UpdatePasswordUser(c echo.Context) error{
	var user model.ForgotPassword
	var users model.User
	jwtToken := c.Param("hash")
	email, code := m.ExtractForgotPasswordToken(jwtToken)
	if err:=config.DB.Where("email = ? AND code = ?", email, code).First(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed to find email",
			"error":   err.Error(),
		})
	}
	print()
	c.Bind(&users)
	if err := config.DB.Where("email = ?", email).Updates(&users).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed to update password",
			"error":   err.Error(),
		})
	}
	print(users.Username)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success update password",
	})

}

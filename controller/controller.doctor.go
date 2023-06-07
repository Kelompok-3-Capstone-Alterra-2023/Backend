package controller

import (
	"capstone/config"
	"capstone/lib/email"
	"capstone/middleware"
	"capstone/model"

	"github.com/labstack/echo/v4"
)

func CreateDoctor(c echo.Context) error {
	var doctor model.Doctor
	var otp model.OTP
	c.Bind(&doctor)
	tempOTP := c.FormValue("otp")
	if tempOTP == ""{
		otp.OTP = email.GenerateOTP()
		otp.DoctorEmail = doctor.Email
		if err:= email.SendEmail("test",doctor.Email,otp.OTP); err!= nil{
			return c.JSON(500, map[string]interface{}{
				"message": "Failed to send OTP",
			})
		}
		err := config.DB.Where("doctor_email=?", doctor.Email).Save(&otp).Error
		if err != nil{
			return c.JSON(500, map[string]interface{}{
				"message": "Failed to save doctor email",
			})	
		}
		return c.JSON(200, map[string]interface{}{
			"message": "Please check your email",
		})
	}else{
		if err:=config.DB.Where("doctor_email = ? AND otp = ?", doctor.Email, tempOTP).First(&otp).Error; err != nil{
			return c.JSON(500, map[string]interface{}{
				"message": "Wrong OTP",
			})
		}
		if err := config.DB.Create(&doctor).Error; err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "Failed to  create doctor",
			})
		}
	}
	return c.JSON(200, map[string]interface{}{
		"message": "success create doctor",
		"data": doctor,
	})
}

func LoginDoctor(c echo.Context) error{
	var doctor model.Doctor
	c.Bind(&doctor)
	if err := config.DB.Where("email = ? AND password = ?", doctor.Email, doctor.Password).First(&doctor).Error; err != nil{
		return c.JSON(500, map[string]interface{}{
			"message": "failed to login",
			"error": err.Error(),
		})
	}
	token, err := middleware.CreateDoctorJWT(doctor.ID)
	if err != nil{
		return c.JSON(500, map[string]interface{}{
			"message": "failed to login",
			"error": err.Error(),
		})
	}
	return c.JSON(200, map[string]interface{}{
		"message": "success login",
		"token": token,
	})
}
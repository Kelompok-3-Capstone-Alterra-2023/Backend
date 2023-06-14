package controller

import (
	"errors"
	"net/http"
	"strconv"

	"capstone/config"
	"capstone/model"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	"capstone/lib/email"
	"capstone/middleware"
)

// for admin
type DoctorAdminController struct{}

// get all doctors
func (a *DoctorAdminController) GetDoctors(c echo.Context) error {
	var doctors []model.Doctor

	config.DB.Find(&doctors)

	if err := config.DB.Find(&doctors).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get all doctors",
		"doctors": doctors,
	})
}

// get doctor by id
func (a *DoctorAdminController) GetDoctor(c echo.Context) error {
	// Bind request data to doctor struct
	doctor := model.Doctor{}
	if err := c.Bind(&doctor); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Find doctor by ID
	id, _ := strconv.Atoi(c.Param("id"))
	if err := config.DB.Find(&doctor, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, doctor)
}

// delete doctor
func (a *DoctorAdminController) DeleteDoctor(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid doctor ID"})
	}

	// check if doctor exists
	doctor := model.Doctor{}
	if err := config.DB.First(&doctor, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Doctor not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve doctor"})
	}

	// delete doctor
	if err := config.DB.Unscoped().Delete(&doctor).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to delete doctor"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Doctor deleted successfully"})
}

// update doctor
func (a *DoctorAdminController) UpdateDoctor(c echo.Context) error {
	data := echo.Map{
		"message": "success update doctor",
	}

	var doctor model.Doctor
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		data["message"] = err.Error()
		return c.JSON(http.StatusBadRequest, data)
	}

	// load doctor from database
	if err := config.DB.First(&doctor, id).Error; err != nil {
		data["message"] = err.Error()
		return c.JSON(http.StatusBadRequest, data)
	}

	// bind updated data to doctor
	if err := c.Bind(&doctor); err != nil {
		data["message"] = err.Error()
		return c.JSON(http.StatusBadRequest, data)
	}

	// update doctor
	if err := config.DB.Save(&doctor).Error; err != nil {
		data["message"] = err.Error()
		return c.JSON(http.StatusBadRequest, data)
	}

	return c.JSON(http.StatusOK, data)
}

// for doctor
type DoctorDoctorController struct{}

// get all doctors
func (d *DoctorDoctorController) GetDoctors(c echo.Context) error {
	var doctors []model.Doctor

	config.DB.Find(&doctors)

	if err := config.DB.Find(&doctors).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get all doctors",
		"doctors": doctors,
	})
}

func CreateDoctor(c echo.Context) error {
	var doctor model.Doctor
	var otp model.DoctorOTP
	c.Bind(&otp)

	if otp.OTP == "" {
		otp.OTP = email.GenerateOTP()
		if err := email.SendEmail("test", otp.Email, otp.OTP); err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "Failed to send OTP",
				"email":   otp.Email,
			})
		}
		err := config.DB.Where("email=?", otp.Email).Save(&otp).Error
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "Failed to save doctor email",
			})
		}
		return c.JSON(200, map[string]interface{}{
			"message": "Please check your email",
		})
	} else {
		if err := config.DB.Where("email= ? AND otp = ?", otp.Email, otp.OTP).First(&otp).Error; err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "Wrong OTP",
			})
		}
		doctor.Email = otp.Email
		doctor.Password = otp.Password
		doctor.FullName = otp.Fullname
		doctor.DisplayName = otp.Displayname
		doctor.Alumnus = otp.Alumnus
		doctor.Workplace = otp.Workplace
		doctor.PracticeAddress = otp.PracticeAddress
		if err := config.DB.Create(&doctor).Error; err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "Failed to  create doctor",
			})
		}
	}
	return c.JSON(200, map[string]interface{}{
		"message": "success create doctor",
		"data":    doctor,
	})
}

func LoginDoctor(c echo.Context) error {
	var doctor model.Doctor
	c.Bind(&doctor)
	if err := config.DB.Where("email = ? AND password = ?", doctor.Email, doctor.Password).First(&doctor).Error; err != nil {
		return c.JSON(500, map[string]interface{}{
			"message": "failed to login",
			"error":   err.Error(),
		})
	}
	token, err := middleware.CreateDoctorJWT(doctor.ID)
	if err != nil {
		return c.JSON(500, map[string]interface{}{
			"message": "failed to login",
			"error":   err.Error(),
		})
	}
	return c.JSON(200, map[string]interface{}{
		"message": "success login",
		"token":   token,
	})
}

// for user
type DoctorUserController struct{}

// get all doctors
func (u *DoctorUserController) GetDoctors(c echo.Context) error {
	var doctors []model.Doctor

	config.DB.Find(&doctors)

	if err := config.DB.Find(&doctors).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get all doctors",
		"doctors": doctors,
	})
}

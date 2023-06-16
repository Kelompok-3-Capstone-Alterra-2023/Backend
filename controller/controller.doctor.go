package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"capstone/config"
	"capstone/model"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	"capstone/middleware"
)

// for admin
type DoctorAdminController struct{}

// Handler untuk menyetujui pendaftaran dokter
func (a *DoctorAdminController) ApproveDoctor(c echo.Context) error {
	var doctor model.Doctor
	c.Bind(&doctor)

	// Cari dokter berdasarkan ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid doctor ID")
	}

	if err := config.DB.First(&doctor, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "doctor not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to retrieve doctor's data")
	}

	// Jika dokter ditemukan
	doctor.Status = "approved"
	if err := config.DB.Save(&doctor).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save changes")
	}

	return c.JSON(http.StatusOK, "doctor registration approved")
}

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
	c.Bind(&doctor)
	doctor.Status = "notapproved"
	if err := config.DB.Create(&doctor).Error; err != nil {
		return c.JSON(500, map[string]interface{}{
			"message": "failed to create doctor",
			"error":   err.Error(),
		})
	}
	return c.JSON(200, map[string]interface{}{
		"message": "success create doctor",
		"doctor":  doctor,
	})
}

func LoginDoctor(c echo.Context) error {
	var doctor model.Doctor
	c.Bind(&doctor)
	if err := config.DB.Where("email = ? AND password = ? AND status = ?", doctor.Email, doctor.Password, "approved").First(&doctor).Error; err != nil {
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

	if err := config.DB.Where("status=?", "approved").Find(&doctors).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get all doctors",
		"doctors": doctors,
	})
}
type DoctorRecipt struct{}

func (u *DoctorRecipt) GetAllDrugs(c echo.Context) error {
	var drugs []model.Drug

	config.DB.Find(&drugs)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get all drugs",
		"recipt":  drugs,
	})
}

func (u *DoctorRecipt) GetDetailRecipt(c echo.Context) error {
	var recipt model.Recipt

	reciptID := c.Param("id")

	config.DB.Model(&model.Recipt{}).Preload("Drugs").Find(&recipt, reciptID).Omit("Doctor")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get recipt",
		"recipt":  recipt,
	})
}

func (u *DoctorRecipt) CreateRecipt(c echo.Context) error {
	var recipt model.Recipt

	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Massage": "json cant empty",
		})
	}

	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	doctorID := middleware.ExtractDocterIdToken(token)

	// if err := c.Bind(doctorID); err != nil {
	// 	return c.JSON(http.StatusOK, "success create recipt")
	// }

	var drugs []model.Drug

	somebyte, _ := json.Marshal(json_map["drugs"])
	errjson := json.Unmarshal(somebyte, &drugs)

	if errjson != nil {
		return c.JSON(http.StatusInternalServerError, "error when in unmarshal")
	}

	recipt.DoctorID = uint(doctorID)
	recipt.Drugs = drugs

	result := config.DB.Create(&recipt)

	if result.RowsAffected < 1 {
		return c.JSON(http.StatusInternalServerError, "error when create recipt")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success create recipt",
	})
}
package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"capstone/config"
	"capstone/lib/email"
	"capstone/model"
	awss3 "capstone/service/aws"
	"capstone/util"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	"capstone/middleware"
)

// for all
type DoctorAllController struct{}

func (u *DoctorAllController) GetDoctors(c echo.Context) error {
	var doctors []model.Doctor
	if err := config.DB.Where("status=?", "approved").Find(&doctors).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success get all doctors",
		"doctors": doctors,
	})
}


func (u *DoctorAllController) GetDoctor(c echo.Context) error {
	var doctor model.Doctor
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := config.DB.Where("id=?", id).Find(&doctor).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success get doctor",
		"doctor":  doctor,
	})
}

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

	doctor.Password, _ = util.GeneratePass(10)
	// Jika dokter ditemukan
	doctor.Status = "approved"
	parsedTime, _ := time.Parse(time.RFC3339, doctor.BirthDate)
	doctor.BirthDate = parsedTime.Format("2006-01-02")
	if err := config.DB.Save(&doctor).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save changes")
	}

	emailContent := fmt.Sprintf("Email: %s\nTemporary Password: %s", doctor.Email, doctor.Password)
	if err := email.SendEmail(doctor.FullName, doctor.Email, "Credential Prevent", emailContent); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed to send email",
			"error":   err.Error(),
		})
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
	var cvurl, ijazahurl, strurl, sipurl, propicurl string
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

	cv, _ := c.FormFile("cv")
	if cv != nil {
		// upload cv
		if err != nil {
			data["message"] = err.Error()
			return c.JSON(http.StatusBadRequest, data)
		}
		date := time.Now().Format("2006-01-02")
		fileext := filepath.Ext(cv.Filename)
		awsObjCV := awss3.CreateObject(date, "cv", fileext, cv)
		cvurl, err = awss3.UploadFileS3(awsObjCV, cv)
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to upload propic",
				"error":   err.Error(),
			})
		}
		doctor.CV = cvurl
	}

	ijazah, _ := c.FormFile("ijazah")
	if ijazah != nil {
		// upload ijazah
		if err != nil {
			data["message"] = err.Error()
			return c.JSON(http.StatusBadRequest, data)
		}
		date := time.Now().Format("2006-01-02")
		fileext := filepath.Ext(ijazah.Filename)
		awsObjIjazah := awss3.CreateObject(date, "ijazah", fileext, ijazah)
		ijazahurl, err = awss3.UploadFileS3(awsObjIjazah, ijazah)
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to upload propic",
				"error":   err.Error(),
			})
		}
		doctor.Ijazah = ijazahurl
	}

	sip, _ := c.FormFile("sip")
	if sip != nil {
		// upload sip
		if err != nil {
			data["message"] = err.Error()
			return c.JSON(http.StatusBadRequest, data)
		}
		date := time.Now().Format("2006-01-02")
		fileext := filepath.Ext(sip.Filename)
		awsObjSip := awss3.CreateObject(date, "sip", fileext, sip)
		sipurl, err = awss3.UploadFileS3(awsObjSip, sip)
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to upload propic",
				"error":   err.Error(),
			})
		}
		doctor.SIP = sipurl
	}

	str, _ := c.FormFile("str")
	if str != nil {
		// upload str
		if err != nil {
			data["message"] = err.Error()
			return c.JSON(http.StatusBadRequest, data)
		}
		date := time.Now().Format("2006-01-02")
		fileext := filepath.Ext(str.Filename)
		awsObjStr := awss3.CreateObject(date, "str", fileext, str)
		strurl, err = awss3.UploadFileS3(awsObjStr, str)
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to upload propic",
				"error":   err.Error(),
			})
		}
		doctor.STR = strurl
	}

	propic, _ := c.FormFile("propic")
	if propic != nil {
		// upload propic
		if err != nil {
			data["message"] = err.Error()
			return c.JSON(http.StatusBadRequest, data)
		}
		date := time.Now().Format("2006-01-02")
		fileext := filepath.Ext(propic.Filename)
		awsObjPropic := awss3.CreateObject(date, "propic", fileext, propic)
		propicurl, err = awss3.UploadFileS3(awsObjPropic, propic)
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to upload propic",
				"error":   err.Error(),
			})
		}
		doctor.Propic = propicurl
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
	var cvurl, ijazahurl, strurl, sipurl, propicurl string
	var awsObjCV, awsObjIjazah, awsObjSip, awsObjStr, awsObjPropic awss3.S3Object
	c.Bind(&doctor)

	cv, err := c.FormFile("cv")
	if cv != nil{
		if err != nil  || filepath.Ext(cv.Filename) == ".pdf"{
			return c.JSON(500, map[string]interface{}{
				"message": "File has to be .pdf",
			})
		}
		date := time.Now().Format("2006-01-02")
		fileext := filepath.Ext(cv.Filename)
		awsObjCV = awss3.CreateObject(date, "cv", fileext, cv)
	}

	ijazah, err := c.FormFile("ijazah")
	if ijazah != nil {
		if err != nil || filepath.Ext(ijazah.Filename) != ".pdf"{
			return c.JSON(500, map[string]interface{}{
				"message": "File has to be .pdf",
			})
		}
		date := time.Now().Format("2006-01-02")
		fileext := filepath.Ext(ijazah.Filename)
		awsObjIjazah = awss3.CreateObject(date, "ijazah", fileext, ijazah)
	}

	str, err := c.FormFile("str")
	if str != nil {
		if err != nil || filepath.Ext(str.Filename) != ".pdf"{
			return c.JSON(500, map[string]interface{}{
				"message": "File has to be .pdf",
			})
		}
		date := time.Now().Format("2006-01-02")
		fileext := filepath.Ext(str.Filename)
		awsObjStr = awss3.CreateObject(date, "str", fileext, str)
	}

	sip, err := c.FormFile("sip")
	if sip != nil {
		if err != nil || filepath.Ext(sip.Filename) != ".pdf"{
			return c.JSON(500, map[string]interface{}{
				"message": "File has to be .pdf",
			})
		}
		date := time.Now().Format("2006-01-02")
		fileext := filepath.Ext(sip.Filename)
		awsObjSip = awss3.CreateObject(date, "sip", fileext, sip)
	}

		propic, err := c.FormFile("propic")
		if propic != nil {
			if err != nil  || (filepath.Ext(cv.Filename) != ".jpg" && filepath.Ext(cv.Filename) != ".png" && filepath.Ext(cv.Filename) != ".jpeg"){
				return c.JSON(500, map[string]interface{}{
					"message": "File has to be .jpg, .png, or .jpeg",
				})
			}
			date := time.Now().Format("2006-01-02")
			fileext := filepath.Ext(propic.Filename)
			awsObjPropic = awss3.CreateObject(date, "propic",fileext, propic)
		propicurl, err = awss3.UploadFileS3(awsObjPropic, propic)
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to upload propic",
				"error":   err.Error(),
			})
		}
		doctor.Propic = propicurl
	}


	cvurl, err = awss3.UploadFileS3(awsObjCV, cv)
	if err != nil {
		return c.JSON(500, map[string]interface{}{
			"message": "failed to upload cv",
			"error":   err.Error(),
		})
	}

	ijazahurl, err = awss3.UploadFileS3(awsObjIjazah, ijazah)
	if err != nil {
		awss3.DeleteObject(awsObjCV)
		return c.JSON(500, map[string]interface{}{
			"message": "failed to upload ijazah",
			"error":   err.Error(),
		})
	}

	strurl, err = awss3.UploadFileS3(awsObjStr, str)
	if err != nil {
		awss3.DeleteObject(awsObjCV, awsObjIjazah)
		return c.JSON(500, map[string]interface{}{
			"message": "failed to upload str",
			"error":   err.Error(),
		})
	}

	sipurl, err = awss3.UploadFileS3(awsObjSip, sip)
	if err != nil {
		awss3.DeleteObject(awsObjCV, awsObjIjazah, awsObjStr)
		return c.JSON(500, map[string]interface{}{
			"message": "failed to upload sip",
			"error":   err.Error(),
		})
	}

	doctor.CV = cvurl
	doctor.Ijazah = ijazahurl
	doctor.STR = strurl
	doctor.SIP = sipurl

	doctor.Status = "notapproved"
	if err := config.DB.Create(&doctor).Error; err != nil {
		awss3.DeleteObject(awsObjCV, awsObjIjazah, awsObjStr, awsObjSip)
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
		if doctor.Password != "admin" {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to login",
				"error":   err.Error(),
			})
		}
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
		"doctor":  doctor,
	})
}

func (d *DoctorDoctorController) UpdateDoctor(c echo.Context) error {
	var cvurl, ijazahurl, strurl, sipurl, propicurl string
	data := echo.Map{
		"message": "success update doctor",
	}

	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	doctorID, err := middleware.ExtractDocterIdToken(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	var doctor model.Doctor
	if err != nil {
		data["message"] = err.Error()
		return c.JSON(http.StatusBadRequest, data)
	}

	// load doctor from database
	if err := config.DB.First(&doctor, int(doctorID)).Error; err != nil {
		data["message"] = err.Error()
		return c.JSON(http.StatusBadRequest, data)
	}

	// bind updated data to doctor
	if err := c.Bind(&doctor); err != nil {
		data["message"] = err.Error()
		return c.JSON(http.StatusBadRequest, data)
	}

	cv, _ := c.FormFile("cv")
	if cv != nil {
		if err != nil || filepath.Ext(cv.Filename) != ".pdf"{
			data["message"] = "File must be .pdf"
			return c.JSON(http.StatusBadRequest, data)
		}
		date := time.Now().Format("2006-01-02")
		fileext := filepath.Ext(cv.Filename)
		awsObjCV := awss3.CreateObject(date, "cv", fileext, cv)
		cvurl, err = awss3.UploadFileS3(awsObjCV, cv)
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to upload propic",
				"error":   err.Error(),
			})
		}
		doctor.CV = cvurl
	}

	ijazah, _ := c.FormFile("ijazah")
	if ijazah != nil {
		// upload ijazah
		if err != nil || filepath.Ext(ijazah.Filename) != ".pdf"{
			data["message"] = "File must be .pdf"
			return c.JSON(http.StatusBadRequest, data)
		}
		fileext := filepath.Ext(ijazah.Filename)
		date := time.Now().Format("2006-01-02")
		awsObjIjazah := awss3.CreateObject(date, "ijazah", fileext, ijazah)
		ijazahurl, err = awss3.UploadFileS3(awsObjIjazah, ijazah)
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to upload propic",
				"error":   err.Error(),
			})
		}
		doctor.Ijazah = ijazahurl
	}

	sip, _ := c.FormFile("sip")
	if sip != nil {
		// upload sip
		if err != nil  || filepath.Ext(sip.Filename) != ".pdf"{
			data["message"] = "File must be .pdf"
			return c.JSON(http.StatusBadRequest, data)
		}
		date := time.Now().Format("2006-01-02")
		fileext := filepath.Ext(sip.Filename)
		awsObjSip := awss3.CreateObject(date, "sip", fileext, sip)
		sipurl, err = awss3.UploadFileS3(awsObjSip, sip)
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to upload propic",
				"error":   err.Error(),
			})
		}
		doctor.SIP = sipurl
	}

	str, _ := c.FormFile("str")
	if str != nil {
		// upload str
		if err != nil  || filepath.Ext(str.Filename) != ".pdf"{
			data["message"] = "File must be .pdf"
			return c.JSON(http.StatusBadRequest, data)
		}
		date := time.Now().Format("2006-01-02")
		fileext := filepath.Ext(str.Filename)
		awsObjStr := awss3.CreateObject(date, "str", fileext, str)
		strurl, err = awss3.UploadFileS3(awsObjStr, str)
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to upload propic",
				"error":   err.Error(),
			})
		}
		doctor.STR = strurl
	}

	propic, _ := c.FormFile("propic")
	if propic != nil {
		if err != nil  || (filepath.Ext(propic.Filename) != ".jpg" && filepath.Ext(propic.Filename) != ".jpeg" && filepath.Ext(propic.Filename) != ".png"){
			data["message"] = "File must be .jpg, .jpeg, or .png"
			return c.JSON(http.StatusBadRequest, data)
		}
		date := time.Now().Format("2006-01-02")
		fileext := filepath.Ext(propic.Filename)
		awsObjPropic := awss3.CreateObject(date, "propic", fileext, propic)
		propicurl, err = awss3.UploadFileS3(awsObjPropic, propic)
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to upload propic",
				"error":   err.Error(),
			})
		}
		doctor.Propic = propicurl
	}
	parsedTime, _ := time.Parse(time.RFC3339, doctor.BirthDate)
	doctor.BirthDate = parsedTime.Format("2006-01-02")
	if err := config.DB.Save(&doctor).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save changes")
	}

	// update doctor
	if err := config.DB.Save(&doctor).Error; err != nil {
		data["message"] = err.Error()
		return c.JSON(http.StatusBadRequest, data)
	}

	return c.JSON(http.StatusOK, data)
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
	for i := range doctors {
		YearIn, _ := strconv.Atoi(doctors[i].YearEntry)
		YearOuts, _ := strconv.Atoi(doctors[i].YearOut)
		doctors[i].WorkExperience = uint(YearOuts - YearIn)
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
	doctorID, err := middleware.ExtractDocterIdToken(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
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
	user_id, errconv := strconv.Atoi(fmt.Sprintf("%v", doctorID))
	if errconv != nil {
		log.Println("error when convert user id in ft create recipt")
	}
	recipt.UserID = uint(user_id)

	result := config.DB.Create(&recipt)

	if result.RowsAffected < 1 {
		return c.JSON(http.StatusInternalServerError, "error when create recipt")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success create recipt",
	})
}

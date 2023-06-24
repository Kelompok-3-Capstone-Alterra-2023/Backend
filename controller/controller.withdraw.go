package controller

import (
	"capstone/middleware"
	"capstone/model"
	"capstone/service/database"
	"capstone/util"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/goodsign/monday"
	"github.com/labstack/echo/v4"
)

type WithdrawController struct{}

func (controller *WithdrawController) RequestWithdraw(c echo.Context) error {
	var withdraw model.Withdraw
	c.Bind(&withdraw)
	if withdraw.Amount>=5000000{
		withdraw.TransactionFee=100000
	}else{
		withdraw.TransactionFee=50000
	}
	withdraw.Total = withdraw.Amount + withdraw.TransactionFee
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	doctorID, err := middleware.ExtractDocterIdToken(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	doctor_ID := strconv.Itoa(int(doctorID))
	doctor, err := database.GetDoctorById(doctor_ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	if doctor.Balance <= withdraw.Total {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "balance tidak cukup",
		})
	}

	withdraw.AccountName = doctor.FullName
	withdraw.Doctor.Email = doctor.Email
	withdraw.DoctorID = uint(doctorID)
	withdraw.Status = "queued"
	withdraw.ReferenceNumber = util.GenerateRandomReferenceNumber()
	err = database.SaveWithdraw(&withdraw)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	date := monday.Format(time.Now(), "02 January 2006", monday.LocaleIdID)
	response := model.WithdrawResponse{
		ReferenceNumber: withdraw.ReferenceNumber,
		Method:          withdraw.Method,
		Bank:            withdraw.Bank,
		AccountName:     withdraw.AccountName,
		AccountNumber:   withdraw.AccountNumber,
		Amount:          withdraw.Amount,
		TransactionFee:  withdraw.TransactionFee,
		Total:           withdraw.Total,
		Date:            date,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    response,
	})
}

func (controller *WithdrawController) GetWithdraws(c echo.Context) error {
	var withdraws []model.WithdrawForGet
	var err error

	if c.QueryParam("keyword") != "" {
		withdraws, err = database.SearchWithdraw(c.QueryParam("keyword"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}
	} else {
		withdraws, err = database.GetWithdraws()
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}
	}

	var response []model.WithdrawsResponse
	for i := range withdraws {
		var withdraw model.WithdrawsResponse
		withdraw.Id = withdraws[i].ID
		withdraw.ReferenceNumber = withdraws[i].ReferenceNumber
		withdraw.Method = withdraws[i].Method
		withdraw.Bank = withdraws[i].Bank
		withdraw.DoctorName = withdraws[i].Doctor.FullName
		withdraw.DoctorEmail = withdraws[i].Doctor.Email
		withdraw.AccountNumber = withdraws[i].AccountNumber
		withdraw.Amount = withdraws[i].Amount
		withdraw.TransactionFee = withdraws[i].TransactionFee
		withdraw.Total = withdraws[i].Total
		withdraw.Date = withdraws[i].CreatedAt.Format("02/01/2006")

		response = append(response, withdraw)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    response,
	})
}

func (controller *WithdrawController) ManageWithdraw(c echo.Context) error {
	var withdraw model.Withdraw
	c.Bind(&withdraw)

	if withdraw.Status == "terima" {
		withdraw, err := database.GetWithdrawByID(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}

		doctor, err := database.GetDoctorById(strconv.Itoa(int(withdraw.DoctorID)))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}

		err = database.UpdateStatusWithdraw(c.Param("id"), "approved")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}

		newBalance := doctor.Balance - withdraw.Total
		err = database.UpdateBalanceDoctor(strconv.Itoa(int(doctor.ID)), newBalance)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}
	} else if withdraw.Status == "tolak" {
		err := database.UpdateStatusWithdraw(c.Param("id"), "rejected")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "success manage withdraw",
	})
}

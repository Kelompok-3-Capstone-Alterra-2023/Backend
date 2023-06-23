package controller

import (
	"capstone/middleware"
	"capstone/model"
	"capstone/service/database"
	"capstone/service/midtrans"
	"capstone/util"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/goodsign/monday"
	"github.com/labstack/echo/v4"
)

type OrderController struct{}

func (controller *OrderController) GetDetailDoctor(c echo.Context) error {
	doctor, err := database.GetDoctorById(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"message": err.Error(),
		})
	}

	yearEntryWork, _ := strconv.Atoi(doctor.YearEntry)
	yearOutWork, _ := strconv.Atoi(doctor.YearOut)
	doctorExperience := uint(yearOutWork - yearEntryWork)
	response := model.OrderDetailDoctorResponse{
		ID:              doctor.ID,
		FullName:        doctor.FullName,
		Propic:          doctor.Propic,
		Specialist:      doctor.Specialist,
		Description:     doctor.Description,
		WorkExperience:  doctorExperience,
		Price:           doctor.Price,
		Alumnus:         doctor.Alumnus,
		PracticeAddress: doctor.PracticeAddress,
		STRNumber:       doctor.STRNumber,
		OnlineStatus:    doctor.StatusOnline,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get detail doctor",
		"data":    response,
	})
}

func (controller *OrderController) Order(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	var booking model.Booking
	var order model.Order
	if err := c.Bind(&booking); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	schedule, err := monday.Parse("Monday, 02 January 2006 15:04:05 MST", booking.Schedule, monday.LocaleIdID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	doctorID := c.Param("id")
	userID := int(middleware.ExtractUserIdToken(token))
	doctor, _ := database.GetDoctorById(doctorID)
	user, _ := database.GetUserById(strconv.Itoa(userID))

	totalAmount := booking.Price + booking.ServiceFee
	var orderNumber string
	for {
		orderNumber = util.GenerateRandomOrderNumber()
		err := database.CheckOrderNumber(orderNumber)
		if err.Error() == "record not found" {
			break
		}
	}

	midtransReq := model.MidtransRequest{
		OrderNumber: orderNumber,
		Amount:      int64(totalAmount),
		Doctor: struct {
			ID       uint
			FullName string
			Price    int64
		}{
			doctor.ID,
			doctor.FullName,
			int64(doctor.Price),
		},
		QTY:        1,
		Method:     booking.ConsultationMethod,
		ServiceFee: int64(booking.ServiceFee),
		User: struct {
			FName string
			Email string
			Phone string
		}{
			user.Username,
			user.Email,
			user.Telp,
		},
	}

	bookingResp, err := midtrans.CreateSnapToken(&midtransReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	order.UserID = uint(userID)
	order.DoctorID = doctor.ID
	order.OrderNumber = orderNumber
	order.Date = time.Now()
	order.SnapToken = bookingResp.Token
	order.PaymentURL = bookingResp.RedirectURL

	err = database.SaveOrder(&order)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	payment := model.Payment{
		OrderID:        order.ID,
		TotalPrice:     totalAmount,
		TransferStatus: "pending",
	}
	err = database.SavePayment(&payment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	consultationSchedule := model.ConsultationSchedule{
		DoctorID: doctor.ID,
		UserID:   uint(userID),
		OrderID:  order.ID,
		Method:   booking.ConsultationMethod,
		Schedule: schedule,
	}
	err = database.SaveSchedule(&consultationSchedule)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success booking",
		"data":    bookingResp,
	})

}

func (controller *OrderController) CheckSchedule(c echo.Context) error {
	schedules, err := database.CheckScheduleByDoctorId(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	for i := range schedules {
		schedule, _ := time.Parse("2006-01-02T15:04:05Z", schedules[i].Schedule)
		schedules[i].Schedule = schedule.Format("2006-01-02 15:04:05")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get schedules",
		"data":    schedules,
	})
}

func (controller *OrderController) GetSchedules(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	doctorID, err := middleware.ExtractDocterIdToken(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	doctor_id := fmt.Sprintf("%f", doctorID)
	schedules, err := database.GetScheduleByDoctor(doctor_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	var response []model.Schedules
	for i := range schedules {
		var schedule model.Schedules
		date := schedules[i].Schedule.Format("02/01/2006 15:04")
		schedule.ID = schedules[i].ID
		schedule.Method = schedules[i].Method
		schedule.Date = date
		schedule.Status = schedules[i].Status
		schedule.UserGender = schedules[i].User.Gender
		schedule.UserName = schedules[i].User.Fullname
		response = append(response, schedule)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    response,
	})
}

func (controller *OrderController) MidtransNotification(c echo.Context) error {
	var notification model.Notification
	if err := c.Bind(&notification); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	if notification.PaymentStatus == "settlement" {
		payment, UserID, doctorID, err := database.GetPaymentandDoctorID(notification.OrderID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}
		doctor_id := strconv.Itoa(int(doctorID))
		doctor, err := database.GetDoctorById(doctor_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}
		newBalance := doctor.Balance + payment.TotalPrice
		err = database.UpdateBalanceDoctor(doctor_id, newBalance)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}

		err = database.UpdateStatusSchedule(strconv.Itoa(int(payment.OrderID)), "menunggu")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}
		user, _ := database.GetUserById(strconv.Itoa(int(UserID)))

		//creating chat room
		_, errroom := createChatRoom(user, doctor)

		if errroom != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "success booking but fail get chatroom",
			})
		}
	}

	err := database.UpdatePayment(&notification)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "success",
	})

}

func (controller *OrderController) OrderManual(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	var booking model.Booking
	var order model.Order
	if err := c.Bind(&booking); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	schedule, err := monday.Parse("Monday, 02 January 2006 15:04:05 MST", booking.Schedule, monday.LocaleIdID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	doctorID := c.Param("id")
	userID := int(middleware.ExtractUserIdToken(token))
	doctor, _ := database.GetDoctorById(doctorID)
	user, _ := database.GetUserById(strconv.Itoa(userID))
	fmt.Println("d", doctorID, userID)

	var orderNumber string
	for {
		orderNumber = util.GenerateRandomOrderNumber()
		err := database.CheckOrderNumber(orderNumber)
		if err.Error() == "record not found" {
			break
		}
	}

	order.UserID = uint(userID)
	order.DoctorID = doctor.ID
	order.OrderNumber = orderNumber
	order.Date = time.Now()

	err = database.SaveOrder(&order)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	totalAmount := booking.Price + booking.ServiceFee
	payment := model.Payment{
		OrderID:         order.ID,
		PaymentMethod:   booking.PaymentMethod,
		TotalPrice:      totalAmount,
		TransferStatus:  "success",
		TransactionTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	err = database.SavePayment(&payment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	newBalance := doctor.Balance + payment.TotalPrice
	err = database.UpdateBalanceDoctor(doctorID, newBalance)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	consultationSchedule := model.ConsultationSchedule{
		DoctorID: doctor.ID,
		UserID:   uint(userID),
		OrderID:  order.ID,
		Method:   booking.ConsultationMethod,
		Status:   "menunggu",
		Schedule: schedule,
	}
	err = database.SaveSchedule(&consultationSchedule)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	//creating chat room
	chatroom, errroom := createChatRoom(user, doctor)

	if errroom != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "success booking but fail get chatroom",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "success booking and get chatroom",
		"chat room id": chatroom,
	})

}

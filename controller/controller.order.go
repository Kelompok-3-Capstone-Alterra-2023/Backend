package controller

import (
	"capstone/model"
	"capstone/service/database"
	"capstone/service/midtrans"
	"capstone/util"
	"net/http"
	"strconv"
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

	yearEntryWork, _ := strconv.Atoi(doctor.DateOfEntry)
	yearOutWork, _ := strconv.Atoi(doctor.DateOfOut)
	doctorExperience := uint(yearOutWork - yearEntryWork)
	response := model.OrderDetailDoctorResponse{
		ID:              doctor.ID,
		FullName:        doctor.FullName,
		Photo:           doctor.Photo,
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

	doctorID := strconv.Itoa(int(booking.DoctorID))
	doctor, _ := database.GetDoctorById(doctorID)
	user, _ := database.GetUserById(c.Param("user_id"))

	totalAmount := booking.Price + booking.ServiceFee
	orderNumber := util.GenerateRandomOrderNumber()
	midtransReq := model.MidtransRequest{
		OrderNumber: orderNumber,
		Amount:      int64(totalAmount),
		Doctor: struct {
			ID       uint
			FullName string
			Price    int64
		}{
			booking.DoctorID,
			doctor.FullName,
			int64(doctor.Price),
		},
		QTY:        1,
		Method:     booking.Method,
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

	userID, _ := strconv.Atoi(c.Param("user_id"))
	order.UserID = uint(userID)
	order.DoctorID = booking.DoctorID
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
		DoctorID: booking.DoctorID,
		UserID:   uint(userID),
		OrderID:  order.ID,
		Method:   booking.Method,
		Status:   "menunggu",
		Schedule: schedule,
	}
	err = database.SaveSchedule(&consultationSchedule)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
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

func (controller *OrderController) Notification(c echo.Context) error {
	var notification model.Notification
	if err := c.Bind(&notification); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
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

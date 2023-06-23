package database

import (
	"capstone/config"
	"capstone/model"
	"errors"
)

func GetDoctorById(id string) (model.Doctor, error) {
	var doctor model.Doctor
	if err := config.DB.First(&doctor, id).Error; err != nil {
		return doctor, err
	}

	return doctor, nil
}

func GetUserById(id string) (model.User, error) {
	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return user, err
	}

	return user, nil
}

func SaveOrder(order *model.Order) error {
	return config.DB.Save(order).Error
}

func SavePayment(payment *model.Payment) error {
	return config.DB.Save(payment).Error
}

func SaveSchedule(schedule *model.ConsultationSchedule) error {
	return config.DB.Save(schedule).Error
}

func CheckScheduleByDoctorId(id string) ([]model.ConsultationScheduleResponse, error) {
	var schedule []model.ConsultationScheduleResponse
	if err := config.DB.Table("consultation_schedules").Where("doctor_id = ? AND schedule >= CURDATE() AND schedule <= DATE_ADD(NOW(), INTERVAL 7 DAY)", id).Scan(&schedule).Error; err != nil {
		return nil, err
	}

	return schedule, nil
}

func UpdatePayment(paymentUpdate *model.Notification) error {
	var payment model.Payment
	var id uint
	if err := config.DB.Table("orders").Select("id").Where("order_number = ?", paymentUpdate.OrderID).Scan(&id).Error; err != nil {
		return err
	}

	if err := config.DB.Model(&payment).Where("order_id = ?", id).Updates(model.Payment{
		OrderID: id, PaymentMethod: paymentUpdate.PaymentType, TransferStatus: paymentUpdate.PaymentStatus, TransactionTime: paymentUpdate.TransactionTime,
	}).Error; err != nil {
		return err
	}

	return nil
}

func CheckOrderNumber(orderNumber string) error {
	var order model.Order
	if err := config.DB.Table("orders").Where("order_number = ?", orderNumber).First(&order).Error; err != nil {
		return err
	}

	return errors.New("")
}

func GetPaymentandDoctorID(orderNumber string) (model.Payment, uint, uint, error) {
	var order model.Order
	var payment model.Payment

	if err := config.DB.Preload("Doctor").Where("order_number = ?", orderNumber).First(&order).Error; err != nil {
		return payment, 0, 0, err
	}
	config.DB.Table("payments").Where("order_id = ?", order.ID).Scan(&payment)

	return payment, order.UserID, order.DoctorID, nil
}

func UpdateStatusSchedule(order_id, status string) error {
	if err := config.DB.Table("consultation_schedules").Where("order_id = ?", order_id).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

func GetScheduleByDoctor(doctor_id string) ([]model.ConsultationSchedule, error) {
	var schedules []model.ConsultationSchedule
	if err := config.DB.Preload("User").Where("doctor_id = ? AND status IS NOT NULL", doctor_id).Find(&schedules).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}

package model

import (
	"time"

	"gorm.io/gorm"
)

type ConsultationSchedule struct {
	gorm.Model
	DoctorID uint      `json:"doctor_id" form:"doctor_id"`
	UserID   uint      `json:"user_id" form:"user_id"`
	OrderID  uint      `json:"order_id" form:"order_id"`
	Method   string    `json:"method" form:"method"`
	Status   string    `json:"status" form:"status"`
	Schedule time.Time `json:"schedule" form:"schedule" gorm:"type:datetime"`
	Doctor   Doctor    `gorm:"foreignKey:DoctorID"`
	User     User      `gorm:"foreignKey:UserID"`
	Order    Order     `gorm:"foreignKey:OrderID"`
}

type ConsultationScheduleResponse struct {
	ID       uint   `json:"id" form:"id"`
	DoctorID uint   `json:"doctor_id" form:"doctor_id"`
	Schedule string `json:"schedule" form:"schedule"`
}
type Schedules struct {
	ID         uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	UserName   string `json:"user_name"`
	UserGender string `json:"user_gender"`
	Method     string `json:"method"`
	Status     string `json:"status"`
	Date       string `json:"date"`
}

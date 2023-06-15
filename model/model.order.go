package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID      uint      `json:"user_id" form:"user_id"`
	DoctorID    uint      `json:"doctor_id" form:"doctor_id"`
	OrderNumber string    `json:"order_number" form:"order_number" gorm:"unique"`
	Date        time.Time `json:"date" form:"date"`
	SnapToken   string    `json:"snap_token" form:"snap_toke"`
	PaymentURL  string    `json:"payment_url" form:"payment_url"`
	User        User      `gorm:"foreignKey:UserID"`
	Doctor      Doctor    `gorm:"foreignKey:DoctorID"`
}

type Booking struct {
	DoctorID   uint    `json:"doctor_id" form:"doctor_id"`
	UserID     uint    `json:"user_id" form:"user_id"`
	Schedule   string  `json:"schedule" form:"schedule"`
	Method     string  `json:"method" form:"method"`
	ServiceFee float64 `json:"service_fee" form:"service_fee"`
	Price      float64 `json:"price" form:"price"`
}

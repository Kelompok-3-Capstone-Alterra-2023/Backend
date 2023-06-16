package model

import (
	"gorm.io/gorm"
)

type DoctorOTP struct {
	gorm.Model
	Email           string `gorm:"unique;not null" json:"email" form:"email"`
	Password        string `gorm:"not null" json:"password" form:"password"`
	Fullname        string `gorm:"not null" json:"fullname" form:"fullname"`
	Displayname     string `json:"displayname" form:"displayname"`
	Alumnus         string ` json:"alumnus" form:"alumnus"`
	Workplace       string `json:"workplace" form:"workplace"`
	PracticeAddress string `json:"practice_address" form:"practice_address"`
	Price           float64
	Balance         float64
	OTP 		   string    `json:"otp" form:"otp"`
}

type UserOTP struct {
	gorm.Model
	Email         string `json:"email" form:"email" gorm:"type:varchar(255)unique;not null"`
	Username      string `json:"username" form:"username" gorm:"type:varchar(255)unique;not null"`
	Password      string `json:"password" form:"password" gorm:"not null"`
	Telp          string `json:"telpon" form:"telpon" gorm:"varchar(20)"`
	Alamat        string `json:"alamat" form:"alamat" gorm:"type:text"`
	Gender        string `json:"gender" form:"gender" gorm:"type:varchar(2)"`
	Status_Online bool   `json:"status_online" form:"status_online" gorm:"type:boolean"`
	OTP 		   string    `json:"otp" form:"otp"`
}


package model

import (
	"gorm.io/gorm"
)

type DoctorOTP struct {
	gorm.Model
	Email           string `gorm:"not null" json:"email" form:"email"`
	Password        string `gorm:"not null" json:"password" form:"password"`
	Fullname        string `gorm:"not null" json:"fullname" form:"fullname"`
	Displayname     string `json:"displayname" form:"displayname" gorm:"type:varchar(255)not null"`
	Alumnus         string ` json:"alumnus" form:"alumnus" gorm:"type:varchar(255)not null"`
	Work       string `json:"workplace" form:"workplace" gorm:"type:varchar(255)not null"`
	PracticeAddress string `json:"practice_address" form:"practice_address" gorm:"type:text not null"`
	Price           float64
	Balance         float64
	OTP 		   string    `json:"otp" form:"otp"`
	Photo 		 string `json:"photo" form:"photo"`
	Status_Online bool   `json:"status_online" form:"status_online" gorm:"type:boolean"`
	Status 	   string `json:"status" form:"status" gorm:"type:varchar(10)"`
	Specialist string `json:"specialist" form:"specialist" gorm:"type:varchar(255)"`
	Description string `json:"description" form:"description" gorm:"type:text"`
}

type UserOTP struct {
	gorm.Model
	Email         string `json:"email" form:"email" gorm:"type:varchar(255);not null"`
	Username      string `json:"username" form:"username" gorm:"type:varchar(255);not null"`
	Password      string `json:"password" form:"password" gorm:"not null"`
	Telp          string `json:"telpon" form:"telpon" gorm:"varchar(20)"`
	Alamat        string `json:"alamat" form:"alamat" gorm:"type:text"`
	Gender        string `json:"gender" form:"gender" gorm:"type:varchar(2)"`
	Status_Online bool   `json:"status_online" form:"status_online" gorm:"type:boolean"`
	OTP 		   string    `json:"otp" form:"otp"`
}


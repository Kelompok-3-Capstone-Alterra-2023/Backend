package model

import (
	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model
	Email           string  `gorm:"unique;not null" json:"email" form:"email"`
	Password        string  `gorm:"not null" json:"password" form:"password"`
	FullName        string  `gorm:"not null" json:"full_name" form:"full_name"`
	DisplayName     string  `json:"display_name" form:"display_name"`
	Alumnus         string  `json:"alumnus" form:"alumnus"`
	Workplace       string  `json:"workplace" form:"workplace"`
	DateOfEntry     string  `json:"date_of_entry" form:"date_of_entry"`
	DateOfOut       string  `json:"date_of_out" form:"date_of_out"`
	PracticeAddress string  `json:"practice_address" form:"practice_address"`
	Price           float64 `json:"price" form:"price" gorm:"type:double"`
	Balance         float64 `json:"balance" form:"balance" gorm:"type:double"`
	Photo           string  `json:"photo" form:"photo"`
	StatusOnline    bool    `json:"status_online" form:"status_online"`
	Status          string  `json:"status" form:"status"`
	STRNumber       string  `json:"str_number" form:"str_number"`
	Specialist      string  `json:"specialist" form:"specialist"`
	Description     string  `json:"description" form:"description"`
}

type OrderDetailDoctorResponse struct {
	ID              uint    `json:"id" form:"id"`
	FullName        string  `json:"full_name" form:"full_name"`
	Photo           string  `json:"photo" form:"phot"`
	Specialist      string  `json:"specialist" form:"specialist"`
	Description     string  `json:"description" form:"description"`
	WorkExperience  uint    `json:"work_experience" form:"work_experience"`
	Price           float64 `json:"price" form:"price"`
	Alumnus         string  `json:"alumnus" form:"alumnus"`
	PracticeAddress string  `json:"practice_address" form:"practice_address"`
	STRNumber       string  `json:"str_number" form:"str_number"`
	OnlineStatus    bool    `json:"status_online" form:"status_online"`
}

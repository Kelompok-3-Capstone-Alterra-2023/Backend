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
	NIK 		   string  `gorm:"not null;unique" json:"nik" form:"nik"`
	Gender string `gorm:"not null" json:"gender" form:"gender"`
	BirthPlace string `gorm:"not null" json:"birth_place" form:"birth_place"`
	BirthDate string `gorm:"not null;type:date" json:"birth_date" form:"birth_date"`
	Religion string `gorm:"not null" json:"religion" form:"religion"`
	Alumnus         string  `gorm:"not null" json:"alumnus" form:"alumnus"`
	Jurusan		 string  `gorm:"not null" json:"jurusan" form:"jurusan"`
	Work      string  `json:"work" form:"work"`
	GradYear string  `gorm:"not null" json:"grad_year" form:"grad_year"`
	Alumnus2		string  `json:"alumnus2" form:"alumnus2"`
	Jurusan2 		string  `json:"jurusan2" form:"jurusan2"`
	GradYear2 		string  `json:"grad_year2" form:"grad_year2"`
	DateOfEntry     string  `json:"date_of_entry" form:"date_of_entry" gorm:"type:date"`
	DateOfOut       string  `json:"date_of_out" form:"date_of_out" gorm:"type:date"`
	PracticeAddress string  `gorm:"not null" json:"practice_address" form:"practice_address"`
	Price           float64 `json:"price" form:"price" gorm:"type:double"`
	Balance         float64 `json:"balance" form:"balance" gorm:"type:double"`
	Photo           string  `json:"photo" form:"photo"`
	CV 			string  `gorm:"not null" json:"cv" form:"cv"`
	Ijazah 			string  `gorm:"not null" json:"ijazah" form:"ijazah"`
	STR 			string  `gorm:"not null" json:"str" form:"str"`
	SIP 			string  `gorm:"not null" json:"sip" form:"sip"`
	StatusOnline    bool    `json:"status_online" form:"status_online"`
	Status          string  `json:"status" form:"status"`
	STRNumber       string  `gorm:"not null" json:"str_number" form:"str_number"`
	Specialist      string  `json:"specialist" form:"specialist"`
	Description     string  `json:"description" form:"description"`
}
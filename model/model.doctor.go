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
	NIK             string  `gorm:"not null;unique" json:"nik" form:"nik"`
	Gender          string  `gorm:"not null" json:"gender" form:"gender"`
	BirthPlace      string  `gorm:"not null" json:"birth_place" form:"birth_place"`
	BirthDate       string  `gorm:"not null;type:date" json:"birth_date" form:"birth_date"`
	Religion        string  `gorm:"not null" json:"religion" form:"religion"`
	Alumnus         string  `gorm:"not null" json:"alumnus" form:"alumnus"`
	Jurusan         string  `gorm:"not null" json:"jurusan" form:"jurusan"`
	WorkExperience  uint    `json:"work_experience" form:"work_experience"`
	Work            string  `json:"work" form:"work"`
	GradYear        string  `gorm:"not null" json:"grad_year" form:"grad_year"`
	Alumnus2        string  `json:"alumnus2" form:"alumnus2"`
	Jurusan2        string  `json:"jurusan2" form:"jurusan2"`
	GradYear2       string  `json:"grad_year2" form:"grad_year2"`
	YearEntry    string  `json:"year_entry" form:"year_entry"`
	YearOut      string `json:"year_out" form:"year_out"`
	PracticeAddress string  `gorm:"not null" json:"practice_address" form:"practice_address"`
	Price           float64 `json:"price" form:"price" gorm:"type:double"`
	Komisi int `json:"komisi" form:"komisi"`
	Balance         float64 `json:"balance" form:"balance" gorm:"type:double"`
	CV              string  `gorm:"not null" json:"cv" form:"cv"`
	Ijazah          string  `gorm:"not null" json:"ijazah" form:"ijazah"`
	STR             string  `gorm:"not null" json:"str" form:"str"`
	SIP             string  `gorm:"not null" json:"sip" form:"sip"`
	Propic		  string  `json:"propic" form:"propic"`
	StatusOnline    bool    `json:"status_online" form:"status_online"`
	Status          string  `json:"status" form:"status"`
	STRNumber       string  `gorm:"not null" json:"str_number" form:"str_number"`
	Specialist      string  `json:"specialist" form:"specialist"`
	Description     string  `json:"description" form:"description"`
	ChatwithUser    []User  `gorm:"many2many:Chatrooms"`
}

type OrderDetailDoctorResponse struct {
	ID              uint    `json:"id" form:"id"`
	FullName        string  `json:"full_name" form:"full_name"`
	Propic		  string  `json:"propic" form:"propic"`
	Specialist      string  `json:"specialist" form:"specialist"`
	Description     string  `json:"description" form:"description"`
	WorkExperience  uint    `json:"work_experience" form:"work_experience"`
	Price           float64 `json:"price" form:"price"`
	Alumnus         string  `json:"alumnus" form:"alumnus"`
	PracticeAddress string  `json:"practice_address" form:"practice_address"`
	STRNumber       string  `json:"str_number" form:"str_number"`
	OnlineStatus    bool    `json:"status_online" form:"status_online"`
}

type OrderDetailAdminHistoryResponse struct {
	Id 			uint    `json:"id" form:"id"`
	DoctorName 	string  `json:"doctor_name" form:"doctor_name"`
	DoctorEmail string  `json:"doctor_email" form:"doctor_email"`
	Komisi 		float64 `json:"komisi" form:"komisi"`
	Tanggal 	string  `json:"tanggal" form:"tanggal"`
	CV 			string  `json:"cv" form:"cv"`
	Ijazah 		string  `json:"ijazah" form:"ijazah"`
	STR 		string  `json:"str" form:"str"`
	SIP 		string  `json:"sip" form:"sip"`
}

type KomisiDoctor struct{
	ID uint `json:"id" form:"id"`
	DoctorID uint `json:"doctor_id" form:"doctor_id"`
	TotalPrice float64 `json:"total_price" form:"total_price"`
	CreatedAt string `json:"created_at" form:"created_at"`
}

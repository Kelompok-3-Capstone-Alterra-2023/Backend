package model

type OTP struct{
	OTP string `gorm:"not null"`
	DoctorEmail string
}
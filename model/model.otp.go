package model

type OTP struct {
	OTP         string `gorm:"not null json:"otp" form:"otp"`
	DoctorEmail string
}

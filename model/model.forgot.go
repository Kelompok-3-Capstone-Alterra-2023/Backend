package model

type ForgotPassword struct{
	Email string `json:"email" form:"email" gorm:"type:varchar(100);not null"`
	Code string `json:"code" form:"code" gorm:"type:varchar(100);not null"`
}
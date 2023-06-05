package model

import "gorm.io/gorm"

type Doctor struct {
	gorm.Model
	Email        string `json:"email" form:"email"`
	FullName     string `json:"fulname" form:"fullname"`
	Display_Name string `json:"display_name" form:"display_name"`
}

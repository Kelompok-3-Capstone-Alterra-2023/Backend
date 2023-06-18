package model

import (
	"gorm.io/gorm"
)

// Struct untuk menyimpan informasi pengguna
type Admin struct {
	gorm.Model
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Email   string `json:"email" form:"email"`
}

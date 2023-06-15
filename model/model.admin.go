package model

import (
	"gorm.io/gorm"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v5"
)

// Struct untuk menyimpan informasi pengguna
type Admin struct {
	gorm.Model
	Username string `gorm:"not null" json:"username" form:"username"`
	Password string `gorm:"not null" json:"password" form:"password"`
}

// Pengguna yang sudah ditentukan
var adminIdentity = Admin{
	Username: "admin",
	Password: "password",
}

// Struct untuk menyimpan data permintaan login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Struct untuk menyimpan data token JWT
type JwtClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
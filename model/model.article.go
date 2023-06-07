package model

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Doctor_ID uint   `json:"doctor_id" form:"doctor_id" gorm:"not null"`
	Title     string `json:"title" form:"title" gorm:"type:varchar(100);not null"`
	Thumbnail string `json:"thumbnail" form:"thumbnail" gorm:"type:varchar(255)"`
	Content   string `json:"content" form:"content" gorm:"type:longtext"`
	Category  string `json:"category" form:"category" gorm:"type:varchar(20)"`
	Status    string `json:"status" form:"status" gorm:"type:varchar(15)"`
	Doctor    Doctor `gorm:"foreignKey:Doctor_ID"`
}

type ArticleResponse struct {
	ID        uint   `json:"id" form:"id"`
	Doctor_ID uint   `json:"doctor_id" form:"doctor_id"`
	Title     string `json:"title" form:"title"`
	Thumbnail string `json:"thumbnail" form:"thumbnail"`
	Content   string `json:"content" form:"content"`
	Category  string `json:"category" form:"category"`
}

type DetailArticleResponse struct {
	ID         uint   `json:"id" form:"id"`
	Updated_At string `json:"date" form:"date"`
	// change to doctor_name
	Doctor_Name string `json:"doctor_name" form:"doctor_name"`
	Title       string `json:"title" form:"title"`
	Thumbnail   string `json:"thumbnail" form:"thumbnail"`
	Content     string `json:"content" form:"content"`
	Category    string `json:"category" form:"category"`
	Est_Read    uint   `json:"est_read" form:"est_read"`
}

type AllArticleResponse struct {
	ID        uint   `json:"id" form:"id"`
	Title     string `json:"title" form:"title"`
	Thumbnail string `json:"thumbnail" form:"thumbnail"`
	Category  string `json:"category" form:"category"`
}

type AllArticleDoctorResponse struct {
	ID        uint   `json:"id" form:"id"`
	Title     string `json:"title" form:"title"`
	Thumbnail string `json:"thumbnail" form:"thumbnail"`
	Category  string `json:"category" form:"category"`
	Status    string `json:"status" form:"status"`
}

type AllArticleAdminResponse struct {
	ID         uint   `json:"id" form:"id"`
	Full_Name  string `json:"doctor_name" form:"doctor_name"`
	Title      string `json:"title" form:"title"`
	Category   string `json:"category" form:"category"`
	Created_At string `json:"date" form:"date"`
}

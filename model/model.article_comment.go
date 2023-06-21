package model

import "gorm.io/gorm"

type ArticleComment struct {
	gorm.Model
	ArticleID uint    `json:"article_id" form:"article_id"`
	FullName  string  `json:"full_name" form:"full_name"`
	Comment   string  `json:"comment" form:"comment" gorm:"type:mediumtext"`
	Article   Article `gorm:"foreignKey:ArticleID"`
}

type ArticleCommentResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Comment  string `json:"comment"`
}

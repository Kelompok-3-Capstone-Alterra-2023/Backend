package database

import (
	"capstone/config"
	"capstone/model"
)

// UNIVERSAL
func DeleteArticle(id string) error {
	var article model.Article
	if err := config.DB.Delete(&article, id).Error; err != nil {
		return err
	}

	return nil
}

func SaveArticleComment(comment *model.ArticleComment) error {
	return config.DB.Save(comment).Error
}

func GetArticleComment(id string) ([]model.ArticleCommentResponse, error) {
	var articleComments []model.ArticleCommentResponse
	if err := config.DB.Table("article_comments").Where("article_id = ?", id).Find(&articleComments).Error; err != nil {
		return nil, err
	}

	return articleComments, nil
}

// DOCTOR
func SaveArticle(article *model.Article) error {
	return config.DB.Save(article).Error
}

func UpdateArticle(updatedArticle *model.Article) error {
	var article model.Article
	article.ID = updatedArticle.ID
	if err := config.DB.Model(&article).Updates(model.Article{Doctor_ID: updatedArticle.Doctor_ID,
		Title: updatedArticle.Title, Thumbnail: updatedArticle.Thumbnail, Content: updatedArticle.Content, Category: updatedArticle.Category,
		Status: updatedArticle.Status}).Error; err != nil {
		return err
	}
	return nil
}

func DoctorGetAllArticles(doctor_id string) ([]model.AllArticleDoctorResponse, error) {
	var articlesDoctor []model.AllArticleDoctorResponse
	if err := config.DB.Table("articles").Where("doctor_id = ? AND deleted_at is null", doctor_id).Find(&articlesDoctor).Error; err != nil {
		return nil, err
	}

	return articlesDoctor, nil
}

func DoctorSearchArticles(doctorID, title string) ([]model.AllArticleDoctorResponse, error) {
	var articles []model.AllArticleDoctorResponse
	if err := config.DB.Table("articles").Where("doctor_id = ? AND title LIKE ? AND deleted_at is null", doctorID, "%"+title+"%").Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

// USER
func GetAllArticles() ([]model.AllArticleResponse, error) {
	var articles []model.AllArticleResponse
	if err := config.DB.Table("articles").Where("status = \"disetujui\" AND deleted_at is null").Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func UserGetArticleById(id string) (model.Article, error) {
	var article model.Article
	if err := config.DB.Preload("Doctor").Where("status = \"disetujui\"").First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

func UserSearchArticles(title string) ([]model.AllArticleResponse, error) {
	var articles []model.AllArticleResponse
	if err := config.DB.Table("articles").Where("status = \"disetujui\" AND title LIKE ? AND deleted_at is null", "%"+title+"%").Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

// ADMIN
func AdminGetArticles() ([]model.AllArticleAdminResponse, error) {
	var articles []model.AllArticleAdminResponse
	if err := config.DB.Table("articles").Select("articles.id, articles.thumbnail ,articles.title, articles.content ,articles.category, articles.created_at, doctors.full_name, articles.status").Joins("left join doctors on articles.doctor_id = doctors.id").Where("articles.deleted_at is null").Scan(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func GetArticleById(id string) (model.Article, error) {
	var article model.Article
	if err := config.DB.Preload("Doctor").First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

func GetArticleByStatus(status string) ([]model.AllArticleAdminResponse, error) {
	var articles []model.AllArticleAdminResponse
	if err := config.DB.Table("articles").Select("articles.id, articles.title, articles.content ,articles.category, articles.created_at, doctors.full_name, articles.status").Joins("left join doctors on articles.doctor_id = doctors.id").Where("articles.status = ? AND articles.deleted_at is null", status).Scan(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func UpdateStatusArticle(id, status string) error {
	var article model.Article
	if err := config.DB.Model(&article).Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

func AdminSearchArticles(title string) ([]model.AllArticleAdminResponse, error) {
	var articles []model.AllArticleAdminResponse
	if err := config.DB.Table("articles").Where("title LIKE ? AND deleted_at is null", "%"+title+"%").Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

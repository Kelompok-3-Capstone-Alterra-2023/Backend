// need AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION, and AWS_S3_BUCKET as env variable
package controller

import (
	"capstone/middleware"
	"capstone/model"
	awss3 "capstone/service/aws"
	"capstone/service/database"
	"capstone/util"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/goodsign/monday"
	"github.com/labstack/echo/v4"
)

// ADMIN
type ArticleAdminController struct{}

func (controller *ArticleAdminController) GetArticles(c echo.Context) error {
	articles, err := database.AdminGetArticles()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	for i := range articles {
		date, _ := time.Parse("2006-01-02T15:04:05.999-07:00", articles[i].Created_At)
		articles[i].Created_At = date.Format("02/01/2006")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get articles",
		"data":    articles,
	})
}

func (controller *ArticleAdminController) GetDetailArticle(c echo.Context) error {
	article, err := database.GetArticleById(c.Param("id"))
	if err != nil {
		if err.Error() == "record not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})

	}

	date := monday.Format(article.UpdatedAt, "02 January 2006", monday.LocaleIdID)
	articleResponse := model.DetailArticleResponse{
		ID:          article.ID,
		Updated_At:  date,
		Doctor_Name: article.Doctor.FullName,
		Title:       article.Title,
		Thumbnail:   article.Thumbnail,
		Content:     article.Content,
		Category:    article.Category,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get detail article",
		"data":    articleResponse,
	})
}

func (controller *ArticleAdminController) AcceptArticle(c echo.Context) error {
	var article model.Article
	c.Bind(&article)

	err := database.UpdateStatusArticle(c.Param("id"), article.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success manage article",
	})
}

func (controller *ArticleAdminController) SearchArticles(c echo.Context) error {
	articles, err := database.AdminSearchArticles(c.QueryParam("keyword"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get articles",
		"data":    articles,
	})
}

func (controller *ArticleAdminController) DeleteArticle(c echo.Context) error {
	article, err := database.GetArticleById(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	parsedURL, err := url.Parse(article.Thumbnail)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
	}

	hosts := util.SplitBy(parsedURL.Host, '.')
	awsObj := awss3.S3Object{
		Bucket: hosts[0],
		Key:    parsedURL.Path,
	}

	err = database.DeleteArticle(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	awss3.DeleteObject(awsObj)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "success delete article",
	})
}

// DOCTOR
type ArticleDoctorController struct {
}

func (controller *ArticleDoctorController) AddArticle(c echo.Context) error {
	var article model.Article
	var imageURI string
	var awsObj awss3.S3Object
	c.Bind(&article)

	image, err := c.FormFile("thumbnail")
	if image != nil {
		if err != nil  || (filepath.Ext(image.Filename) != ".jpg" && filepath.Ext(image.Filename) != ".png" && filepath.Ext(image.Filename) != ".jpeg"){
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
		}
		date := time.Now().Format("2006-01-02")
		fileext := filepath.Ext(image.Filename)
		awsObj = awss3.CreateObject(date, "article", fileext, image)
	}

	// get doctor id from jwt token
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	doctorID, err := middleware.ExtractDocterIdToken(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	if awsObj.Key != "" {
		imageURI, err = awss3.UploadFileS3(awsObj, image)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}
	}

	article.Doctor_ID = uint(doctorID)
	article.Thumbnail = imageURI
	article.Status = "ditinjau"
	err = database.SaveArticle(&article)
	if err != nil {
		awss3.DeleteObject(awsObj)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success add article",
	})
}

func (controller *ArticleDoctorController) UpdateArticle(c echo.Context) error {
	var updatedArticle model.Article
	var awsObj awss3.S3Object
	c.Bind(&updatedArticle)

	article, err := database.GetArticleById(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	imageURI := article.Thumbnail
	image, _ := c.FormFile("thumbnail")
	if filepath.Ext(image.Filename) != ".jpg" && filepath.Ext(image.Filename) != ".png" && filepath.Ext(image.Filename) != ".jpeg"{
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "file must be image",
		})
	}
	if image != nil {
		date := time.Now().Format("2006-01-02")
		fileext := filepath.Ext(image.Filename)
		awsObj = awss3.CreateObject(date, "article", fileext, image)
	}

	// get doctor id from jwt token
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	doctorID, err := middleware.ExtractDocterIdToken(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	if awsObj.Key != "" {
		imageURI, err = awss3.UploadFileS3(awsObj, image)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}
	}

	if article.Doctor_ID == uint(doctorID) {
		articleID, _ := strconv.Atoi(c.Param("id"))
		updatedArticle.ID = uint(articleID)
		updatedArticle.Thumbnail = imageURI
		err = database.UpdateArticle(&updatedArticle)
		if err != nil {
			awss3.DeleteObject(awsObj)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}

		articleResponse := model.ArticleResponse{
			ID:        updatedArticle.ID,
			Doctor_ID: uint(doctorID),
			Title:     updatedArticle.Title,
			Thumbnail: updatedArticle.Thumbnail,
			Content:   updatedArticle.Content,
			Category:  updatedArticle.Category,
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success update article",
			"data":    articleResponse,
		})

	}
	return c.JSON(http.StatusUnauthorized, map[string]interface{}{
		"message": "unauthorized",
	})

}

func (controller *ArticleDoctorController) DeleteArticle(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	doctorID, err := middleware.ExtractDocterIdToken(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	doctor_id := uint(doctorID)
	article, err := database.GetArticleById(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	parsedURL, err := url.Parse(article.Thumbnail)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
	}

	hosts := util.SplitBy(parsedURL.Host, '.')
	awsObj := awss3.S3Object{
		Bucket: hosts[0],
		Key:    parsedURL.Path,
	}

	if article.Doctor_ID == doctor_id {
		err = database.DeleteArticle(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}
		awss3.DeleteObject(awsObj)
		return c.JSON(http.StatusOK, map[string]string{
			"message": "success delete article",
		})
	}

	return c.JSON(http.StatusUnauthorized, map[string]string{
		"message": "unauthorized",
	})

}

func (controller *ArticleDoctorController) GetArticles(c echo.Context) error {
	// change with doctor id from jwt
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	doctorID, err := middleware.ExtractDocterIdToken(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	doctor_id := fmt.Sprintf("%f", doctorID)
	articlesDoctor, err := database.DoctorGetAllArticles(doctor_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get articles",
		"data":    articlesDoctor,
	})
}

func (controller *ArticleDoctorController) GetArticle(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	doctorID, err := middleware.ExtractDocterIdToken(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	doctor_id := uint(doctorID)

	article, err := database.GetArticleById(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	if article.Doctor_ID == doctor_id {
		total_word := len(strings.Fields(article.Content))
		est_read := math.Round(float64(total_word) / 200.0)
		date := monday.Format(article.UpdatedAt, "02 January 2006", monday.LocaleIdID)
		articleResponse := model.DetailArticleResponse{
			ID:          article.ID,
			Updated_At:  date,
			Doctor_Name: article.Doctor.FullName,
			Title:       article.Title,
			Thumbnail:   article.Thumbnail,
			Content:     article.Content,
			Category:    article.Category,
			Est_Read:    uint(est_read),
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success get article",
			"data":    articleResponse,
		})
	}

	return c.JSON(http.StatusUnauthorized, map[string]string{
		"message": "unauthorized",
	})
}

func (controller *ArticleDoctorController) SearchArticles(c echo.Context) error {
	token := strings.Fields(c.Request().Header.Values("Authorization")[0])[1]
	doctorID, err := middleware.ExtractDocterIdToken(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	doctor_id := fmt.Sprintf("%f", doctorID)
	articles, err := database.DoctorSearchArticles(doctor_id, c.QueryParam("keyword"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get articles",
		"data":    articles,
	})
}

// USER
type ArticleUserController struct{}

func (controller *ArticleUserController) GetArticles(c echo.Context) error {
	articles, err := database.GetAllArticles()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	if c.QueryParam("category") != "" {
		filteredArticles := []model.AllArticleResponse{}

		for _, article := range articles {
			if article.Category == c.QueryParam("category") {
				filteredArticles = append(filteredArticles, article)
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "succes get articles",
			"data":    filteredArticles,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "succes get articles",
		"data":    articles,
	})
}

func (controller *ArticleUserController) GetDetailArticle(c echo.Context) error {
	article, err := database.UserGetArticleById(c.Param("id"))
	if err != nil {
		if err.Error() == "record not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	articleComment, err := database.GetArticleComment(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	total_word := len(strings.Fields(article.Content))
	est_read := math.Round(float64(total_word) / 200.0)
	date := monday.Format(article.UpdatedAt, "02 January 2006", monday.LocaleIdID)
	articleResponse := model.DetailArticleResponse{
		ID:          article.ID,
		Updated_At:  date,
		Doctor_Name: article.Doctor.FullName,
		Title:       article.Title,
		Thumbnail:   article.Thumbnail,
		Content:     article.Content,
		Category:    article.Category,
		Est_Read:    uint(est_read),
		Comments:    articleComment,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get detail article",
		"data":    articleResponse,
	})
}

func (controller *ArticleUserController) SearchArticles(c echo.Context) error {
	articles, err := database.UserSearchArticles(c.QueryParam("keyword"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	if c.QueryParam("category") != "" {
		filteredArticles := []model.AllArticleResponse{}

		for _, article := range articles {
			if article.Category == c.QueryParam("category") {
				filteredArticles = append(filteredArticles, article)
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "succes get articles",
			"data":    filteredArticles,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get articles",
		"data":    articles,
	})
}

func AddComment(c echo.Context) error {
	var comment model.ArticleComment
	c.Bind(&comment)
	articleID, _ := strconv.Atoi(c.Param("id"))
	comment.ArticleID = uint(articleID)
	err := database.SaveArticleComment(&comment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "success add comment",
	})
}

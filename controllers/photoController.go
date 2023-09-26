package controllers

import (
	"net/http"
	"strconv"

	"final-project-pbi-btpns/app"
	"final-project-pbi-btpns/database"
	"final-project-pbi-btpns/helpers"
	"final-project-pbi-btpns/models"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

func GetPhotoList(c *gin.Context) {
	var photos []app.PhotoResult
	tokenString := c.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// menampilkan foto berdasarkan user_id dan diurutkan berdasarkan created_at  db.Preload("UserResult").
	if err := database.DB.Table("photos").Select("photos.id, photos.title, photos.caption, photos.photo_url, photos.created_at, photos.updated_at, users.email").Joins("JOIN users ON users.id = photos.user_id").Where("photos.user_id = ?", claims.ID).Order("photos.created_at desc").Scan(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": photos})
}

func GetPhotoByID(c *gin.Context) {
	var photo app.PhotoResult
	id := c.Param("photoId")
	tokenString := c.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Table("photos").Select("photos.id, photos.title, photos.caption, photos.photo_url, photos.created_at, photos.updated_at, users.email").Joins("JOIN users ON users.id = photos.user_id").Where("photos.id = ? AND photos.user_id = ?", id, claims.ID).Scan(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if photo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "photo not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": photo})

}

func CreatePhoto(c *gin.Context) {
	var photoCreate app.PhotoCreate
	if err := c.Bind(&photoCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validasi data pengguna
	if _, err := govalidator.ValidateStruct(photoCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	desiredExtensions := []string{"jpg", "jpeg", "png", "gif"}

	// periksa apakah URL foto valid dan berakhir dengan salah satu ekstensi yang diinginkan
	if !helpers.IsValidURLWithDesiredExtension(photoCreate.PhotoUrl, desiredExtensions) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid photo URL or doesn't end with the desired extension"})
		return
	}

	// ambil user_id dari token
	tokenString := c.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// buat objek photo baru
	photo := models.Photo{
		Title:    photoCreate.Title,
		Caption:  photoCreate.Caption,
		PhotoUrl: photoCreate.PhotoUrl,
		UserID:   claims.ID,
	}

	// simpan objek photo ke database
	if err := database.DB.Create(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "photos added successfully"})
}

func UpdatePhoto(c *gin.Context) {
	photoID, err := strconv.Atoi(c.Param("photoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id photo"})
		return
	}
	var photoUpdate app.PhotoUpdate
	if err := c.ShouldBindJSON(&photoUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validasi data pengguna menggunakan govalidator
	if _, err := govalidator.ValidateStruct(photoUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	desiredExtensions := []string{"jpg", "jpeg", "png", "gif"}

	// periksa apakah URL foto valid dan berakhir dengan salah satu ekstensi yang diinginkan
	if !helpers.IsValidURLWithDesiredExtension(photoUpdate.PhotoUrl, desiredExtensions) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid photo URL or doesn't end with the desired extension"})
		return
	}

	// ambil user_id dari token
	tokenString := c.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// Ambil data photo dari database
	var photo models.Photo

	if err := database.DB.Where("id = ? AND user_id = ?", photoID, claims.ID).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "photo not found"})
		return
	}

	photo.Title = photoUpdate.Title
	photo.Caption = photoUpdate.Caption
	photo.PhotoUrl = photoUpdate.PhotoUrl
	if err := database.DB.Save(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "photo updated successfully"})

}

func DeletePhoto(c *gin.Context) {
	var photo models.Photo
	photoID, err := strconv.Atoi(c.Param("photoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id photo"})
		return
	}

	// Ambil user_id dari token
	tokenString := c.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if err := database.DB.First(&photo, photoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "photo not found"})
		return
	}

	if photo.UserID != claims.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you don't have access to delete this photo"})
		return
	}

	if err := database.DB.Delete(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "photo deleted successfully"})
}

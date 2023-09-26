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

// membuat akun pengguna baru
func Register(c *gin.Context) {
	var userRegister app.UserRegister
	if err := c.Bind(&userRegister); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error read userRegister struct"})
		return
	}

	// memvalidasi data pengguna menggunakan govalidator
	if _, err := govalidator.ValidateStruct(userRegister); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
	}

	// cek apakah email sudah terdaftar
	var user models.User

	if len(userRegister.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password at least 6 characters long"})
		return
	}

	if err := database.DB.Where("email = ?", userRegister.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
		return
	}

	// cek apakah username sudah terdaftar
	if err := database.DB.Where("username = ?", userRegister.Username).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already registered"})
		return
	}

	user = models.User{
		Username: userRegister.Username,
		Email:    userRegister.Email,
		Password: userRegister.Password,
	}

	if err := user.HashPassword(user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	record := database.DB.Create(&user)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "successfully created an account"})
}

// login atau autentikasi pengguna
func Login(c *gin.Context) {
	var userLogin app.UserLogin
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "erorr read userLogin struct"})
		return
	}

	// Validasi data pengguna menggunakan govalidator
	if _, err := govalidator.ValidateStruct(userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", userLogin.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	if err := user.CheckPassword(userLogin.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := helpers.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create a token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login succesfully", "your token": token})
}

// mengambil data user berdasarkan id
func GetUserByID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	tokenString := c.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	if userID != int(claims.ID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not allowed"})
		c.Abort()
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", claims.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var userResult app.UserResult
	userResult.ID = user.ID
	userResult.Username = user.Username
	userResult.Email = user.Email
	userResult.CreatedAt = user.CreatedAt.String()
	userResult.UpdatedAt = user.UpdatedAt.String()

	c.JSON(http.StatusOK, gin.H{"data": userResult})
}

// update data user berdasarkan id
func UpdateUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ivalid user id"})
		return
	}
	var userUpdate app.UserUpdate
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validasi data user
	if _, err := govalidator.ValidateStruct(userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	//   password minimal 6 karakter
	if len(userUpdate.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters long"})
		return
	}

	//   validasi email dan username sudah terdaftar atau belum selain data user yang sedang login
	if err := database.DB.Where("email = ? AND id != ?", userUpdate.Email, userID).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
		return
	}

	if err := database.DB.Where("username = ? AND id != ?", userUpdate.Username, userID).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already registered"})
		return
	}

	tokenString := c.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if userID != int(claims.ID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not allowed"})
		return
	}

	if err := database.DB.Where("id = ?", claims.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Username = userUpdate.Username
	user.Email = userUpdate.Email
	if userUpdate.Password != "" {
		if err := user.HashPassword(userUpdate.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully updated user"})
}

// menghapus data user berdasarkan id
func DeleteUser(c *gin.Context) {
	var user models.User
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	tokenString := c.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if userID != int(claims.ID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not allowed"})
		return
	}

	if err := database.DB.Where("id = ?", claims.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted user"})
}

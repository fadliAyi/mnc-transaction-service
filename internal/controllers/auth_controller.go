package controllers

import (
	"net/http"
	"transaction-service/internal/database"
	"transaction-service/internal/models"

	auth "transaction-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Register(c *gin.Context) {
	var input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
		Pin         string `json:"pin"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		UserID:      uuid.New(),
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		PhoneNumber: input.PhoneNumber,
		Address:     input.Address,
		Pin:         input.Pin,
	}

	database.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": user})
}

func Login(c *gin.Context) {
	var input struct {
		PhoneNumber string `json:"phone_number"`
		Pin         string `json:"pin"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check phone number and PIN
	var user models.User
	if err := database.DB.Where("phone_number = ? AND pin = ?", input.PhoneNumber, input.Pin).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Phone Number and PIN doesn’t match."})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.UserID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": gin.H{"access_token": token}})
}

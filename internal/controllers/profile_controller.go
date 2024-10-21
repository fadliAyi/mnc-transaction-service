package controllers

import (
	"net/http"
	"time"
	"transaction-service/internal/database"
	"transaction-service/internal/models"

	"github.com/gin-gonic/gin"
)

type UpdateProfileInput struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Address   string `json:"address" binding:"required"`
}

func UpdateProfile(c *gin.Context) {
	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the userID from the JWT middleware
	userID, _ := c.Get("userID")

	// Find the user by userID
	var user models.User
	if err := database.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update the user's profile fields
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Address = input.Address
	user.UpdatedAt = time.Now() // Update the timestamp for the last modification

	// Save the updated user profile to the database
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Respond with the updated profile
	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"user_id":    user.UserID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"address":    user.Address,
			"updated_at": user.UpdatedAt,
		},
	})
}

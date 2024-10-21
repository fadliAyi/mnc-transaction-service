package controllers

import (
	"net/http"
	"strconv"
	"time"
	"transaction-service/internal/background"
	"transaction-service/internal/database"
	"transaction-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TopUpInput struct {
	Amount int `json:"amount" binding:"required,min=1"`
}

type TransferInput struct {
	TargetUserID string `json:"target_user" binding:"required"`
	Amount       int    `json:"amount" binding:"required,min=1"`
	Remarks      string `json:"remarks"`
}

type PaymentInput struct {
	Amount  int    `json:"amount" binding:"required,min=1"`
	Remarks string `json:"remarks"`
}

func Topup(c *gin.Context) {
	var input TopUpInput
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

	// Calculate the new balance
	newBalance := user.Balance + input.Amount

	// Create a new top-up transaction
	topUpTransaction := models.Transaction{
		TransactionID: uuid.New(),
		UserID:        user.UserID,
		Type:          "CREDIT",
		Amount:        input.Amount,
		Remarks:       "Top-Up",
		BalanceBefore: user.Balance,
		BalanceAfter:  newBalance,
		CreatedAt:     time.Now(),
	}

	// Start a database transaction to ensure atomicity
	tx := database.DB.Begin()
	if err := tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process transaction"})
		return
	}

	// Update the user's balance
	user.Balance = newBalance
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	// Save the transaction
	if err := tx.Create(&topUpTransaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save transaction"})
		return
	}

	// Commit the transaction
	tx.Commit()

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"top_up_id":      topUpTransaction.TransactionID,
			"amount_top_up":  input.Amount,
			"balance_before": topUpTransaction.BalanceBefore,
			"balance_after":  topUpTransaction.BalanceAfter,
			"created_date":   topUpTransaction.CreatedAt,
		},
	})
}

func Transfer(c *gin.Context) {
	var input TransferInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the userID from the JWT middleware
	userID, _ := c.Get("userID")

	// Check if the target user exists
	var targetUser models.User
	if err := database.DB.Where("user_id = ?", input.TargetUserID).First(&targetUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target user not found"})
		return
	}

	// Enqueue the transfer job
	job := background.TransferJob{
		FromUserID: userID.(string),
		ToUserID:   input.TargetUserID,
		Amount:     input.Amount,
		Remarks:    input.Remarks,
	}
	background.TransferQueue <- job // Enqueue the job for background processing

	c.JSON(http.StatusOK, gin.H{
		"status":  "SUCCESS",
		"message": "Transfer is being processed in the background",
	})
}

func Payment(c *gin.Context) {
	var input PaymentInput
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

	// Check if the user has enough balance for the payment
	if user.Balance < input.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Balance is not enough"})
		return
	}

	// Calculate the new balance after the payment
	newBalance := user.Balance - input.Amount

	// Create a new payment transaction
	paymentTransaction := models.Transaction{
		TransactionID: uuid.New(),
		UserID:        user.UserID,
		Type:          "DEBIT",
		Amount:        input.Amount,
		Remarks:       input.Remarks,
		BalanceBefore: user.Balance,
		BalanceAfter:  newBalance,
		CreatedAt:     time.Now(),
	}

	// Start a database transaction to ensure atomicity
	tx := database.DB.Begin()
	if err := tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process transaction"})
		return
	}

	// Update the user's balance
	user.Balance = newBalance
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	// Save the transaction
	if err := tx.Create(&paymentTransaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save transaction"})
		return
	}

	// Commit the transaction
	tx.Commit()

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"payment_id":     paymentTransaction.TransactionID,
			"amount":         input.Amount,
			"balance_before": paymentTransaction.BalanceBefore,
			"balance_after":  paymentTransaction.BalanceAfter,
			"created_date":   paymentTransaction.CreatedAt,
		},
	})
}

func TransactionReport(c *gin.Context) {
	// Get the userID from the JWT middleware
	userID, _ := c.Get("userID")

	// Get pagination parameters
	limit := 10
	offset := 0
	if c.Query("limit") != "" {
		limit, _ = strconv.Atoi(c.Query("limit"))
	}
	if c.Query("offset") != "" {
		offset, _ = strconv.Atoi(c.Query("offset"))
	}

	// Fetch paginated transactions
	var transactions []models.Transaction
	if err := database.DB.Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	// Return the transactions as the response
	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": transactions,
	})
}

package background

import (
	"log"
	"time"
	"transaction-service/internal/database"
	"transaction-service/internal/models"

	"github.com/google/uuid"
)

type TransferJob struct {
	FromUserID string
	ToUserID   string
	Amount     int
	Remarks    string
}

var TransferQueue chan TransferJob

func InitTransferQueue() {
	TransferQueue = make(chan TransferJob, 100) // Buffer size of 100
	go processQueue()
}

func processQueue() {
	for {
		select {
		case job := <-TransferQueue:
			ProcessTransfer(job)
		}
	}
}

func ProcessTransfer(job TransferJob) {
	log.Printf("Processing transfer from user %s to user %s of amount %d", job.FromUserID, job.ToUserID, job.Amount)

	// Start a new database transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		log.Println("Failed to start database transaction:", tx.Error)
		return
	}

	// Find the sender and receiver users
	var fromUser, toUser models.User
	if err := tx.Where("user_id = ?", job.FromUserID).First(&fromUser).Error; err != nil {
		log.Println("Sender not found:", err)
		tx.Rollback()
		return
	}

	if err := tx.Where("user_id = ?", job.ToUserID).First(&toUser).Error; err != nil {
		log.Println("Receiver not found:", err)
		tx.Rollback()
		return
	}

	// Check if the sender has enough balance
	if fromUser.Balance < job.Amount {
		log.Println("Sender does not have enough balance")
		tx.Rollback()
		return
	}

	// Update the balances
	fromUser.Balance -= job.Amount
	toUser.Balance += job.Amount

	// Create the transfer record
	transferID := uuid.New()
	transfer := models.Transaction{
		TransactionID: transferID,
		UserID:        fromUser.UserID,
		Type:          "DEBIT",
		Amount:        job.Amount,
		Remarks:       job.Remarks,
		BalanceBefore: fromUser.Balance + job.Amount, // balance before debit
		BalanceAfter:  fromUser.Balance,              // balance after debit
		CreatedAt:     time.Now(),
	}

	// Create the transfer record for the receiver
	transferReceiver := models.Transaction{
		TransactionID: uuid.New(),
		UserID:        toUser.UserID,
		Type:          "CREDIT",
		Amount:        job.Amount,
		Remarks:       job.Remarks,
		BalanceBefore: toUser.Balance - job.Amount, // balance before credit
		BalanceAfter:  toUser.Balance,              // balance after credit
		CreatedAt:     time.Now(),
	}

	// Save the updated balances and transfer record in the transaction
	if err := tx.Save(&fromUser).Error; err != nil {
		log.Println("Failed to update sender balance:", err)
		tx.Rollback()
		return
	}

	if err := tx.Save(&toUser).Error; err != nil {
		log.Println("Failed to update receiver balance:", err)
		tx.Rollback()
		return
	}

	if err := tx.Create(&transfer).Error; err != nil {
		log.Println("Failed to create transfer record:", err)
		tx.Rollback()
		return
	}

	if err := tx.Create(&transferReceiver).Error; err != nil {
		log.Println("Failed to create transfer record for the receiver:", err)
		tx.Rollback()
		return
	}

	tx.Commit()
	log.Printf("Transfer complete: %s", transferID)
}

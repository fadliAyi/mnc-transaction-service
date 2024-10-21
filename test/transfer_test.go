package test

import (
	"testing"
	"time"
	"transaction-service/internal/background"
	"transaction-service/internal/database"
	"transaction-service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransferWorker(t *testing.T) {
	// Create mock repository
	mockDB := new(database.MockRepository)

	// Initialize the transfer queue with the mock repository
	background.InitTransferQueue()

	// Mock data
	fromUserID := uuid.New().String()
	toUserID := uuid.New().String()
	amount := 50000

	fromUser := &models.User{
		UserID:  uuid.MustParse(fromUserID),
		Balance: 100000,
	}

	toUser := &models.User{
		UserID:  uuid.MustParse(toUserID),
		Balance: 50000,
	}

	// Set up mock expectations
	mockDB.On("Begin").Return(mockDB)
	mockDB.On("Where", "user_id = ?", []interface{}{fromUserID}).Return(mockDB)
	mockDB.On("First", &fromUser).Return(nil)
	mockDB.On("Where", "user_id = ?", []interface{}{toUserID}).Return(mockDB)
	mockDB.On("First", &toUser).Return(nil)
	mockDB.On("Save", &fromUser).Return(nil)
	mockDB.On("Save", &toUser).Return(nil)
	mockDB.On("Create", mock.AnythingOfType("*models.Transaction")).Return(nil)
	mockDB.On("Commit").Return(nil)

	// Create a transfer job
	job := background.TransferJob{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
		Remarks:    "Test transfer",
	}

	// Add the job to the transfer queue
	background.TransferQueue <- job

	// Allow some time for the worker to process the job
	time.Sleep(1 * time.Second)

	// Assert that the mock expectations were met
	mockDB.AssertExpectations(t)

	// Assert balances are updated correctly
	assert.Equal(t, 50000, fromUser.Balance)
	assert.Equal(t, 100000, toUser.Balance)
}

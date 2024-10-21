package database

import (
	"transaction-service/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FindUserByID(userID string) (*models.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) UpdateUserBalance(userID string, newBalance int) error {
	args := m.Called(userID, newBalance)
	return args.Error(0)
}

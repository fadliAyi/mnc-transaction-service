package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	TransactionID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID
	Type          string // DEBIT or CREDIT
	Amount        int
	Remarks       string
	BalanceBefore int
	BalanceAfter  int
	CreatedAt     time.Time
}

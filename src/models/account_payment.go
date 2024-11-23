package models

import (
	"gorm.io/gorm"
	"time"
)

type AccountPayment struct {
	gorm.Model
	ID              uint       `gorm:"primaryKey;autoIncrement"`
	DueDate         time.Time  `json:"due_date"`
	DueAmount       int        `json:"due_amount"`
	PrincipalAmount int        `json:"principal_amount"`
	InterestAmount  int        `json:"interest_amount"`
	PaidDate        *time.Time `json:"paid_date"`
	PaidAmount      *int       `json:"paid_amount"`
	Status          *int       `json:"status"`
	LoanID          uint       `json:"loan_id"`
}

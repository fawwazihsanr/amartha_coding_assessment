package models

import "gorm.io/gorm"

type Loan struct {
	gorm.Model
	ID             uint `gorm:"primaryKey;autoIncrement" json:"id"`
	LoanAmount     int  `json:"loan_amount"`
	InterestAmount int  `json:"interest_amount"`
	LoanDuration   int  `json:"loan_duration"`
}

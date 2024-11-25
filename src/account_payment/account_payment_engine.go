package account_payment_engine

import (
	"amartha_coding/src/models"
	"errors"
	"gorm.io/gorm"
	"time"
)

func GetOutstanding(db *gorm.DB, loanID uint) (int, error) {
	var payments []models.AccountPayment
	var loan models.Loan

	err := db.First(&loan, loanID).Error
	if err != nil {
		return 0, err
	}

	err = db.Where("loan_id = ? AND due_amount > 0", loanID).Find(&payments).Error
	if err != nil {
		return 0, err
	}

	outstanding := 0
	for _, payment := range payments {
		outstanding += payment.DueAmount
	}

	return outstanding, nil
}

func IsDelinquent(db *gorm.DB, loanID uint) (bool, error) {
	var payments []models.AccountPayment

	err := db.Where("loan_id = ?", loanID).Order("due_date ASC").Find(&payments).Error
	if err != nil {
		return false, err
	}

	latePayment := 0
	for _, payment := range payments {
		if time.Now().After(payment.DueDate) && payment.DueAmount > 0 {
			latePayment++
			if latePayment >= 2 {
				return true, nil
			}
		} else {
			latePayment = 0
		}
	}

	return false, nil
}

func MakePayment(db *gorm.DB, loanID uint, paymentAmount int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var payments []models.AccountPayment

		err := tx.Where("loan_id = ? AND due_amount > 0", loanID).Order("due_date ASC").Find(&payments).Error
		if err != nil {
			return errors.New("failed to fetch payments: " + err.Error())
		}

		if len(payments) == 0 {
			return errors.New("no pending payments for this loan")
		}

		remainingAmount := paymentAmount
		paidDate := time.Now()

		for i := range payments {
			if remainingAmount <= 0 {
				break
			}

			payment := &payments[i]
			if remainingAmount >= payment.DueAmount {
				remainingAmount -= payment.DueAmount
				payment.DueAmount = 0
				payment.PrincipalAmount = 0
				payment.InterestAmount = 0
				payment.PaidDate = &paidDate
				payment.PaidAmount = &remainingAmount
			} else {
				payment.DueAmount -= remainingAmount

				if payment.PrincipalAmount > 0 {
					principalReduction := (payment.PrincipalAmount * remainingAmount) / (payment.PrincipalAmount + payment.InterestAmount)
					payment.PrincipalAmount -= principalReduction
					payment.InterestAmount -= remainingAmount - principalReduction
				} else {
					payment.InterestAmount -= remainingAmount
				}

				remainingAmount = 0
			}

			err = tx.Save(payment).Error
			if err != nil {
				return errors.New("failed to update payment record: " + err.Error())
			}
		}

		// Handle any cashback or refund logic here
		if remainingAmount > 0 {

		}

		return nil
	})
}

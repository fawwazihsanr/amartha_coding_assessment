package loan_engine

import (
	"amartha_coding/src/models"
	"gorm.io/gorm"
	"time"
)

func NewLoan(db *gorm.DB, loanAmount int, loanDuration int, firstPaymentDate time.Time) (*models.Loan, error) {
	var createdLoan *models.Loan

	err := db.Transaction(func(tx *gorm.DB) error {
		interestRate := 0.1
		totalInterest := int(float64(loanAmount) * interestRate)

		createdLoan = &models.Loan{
			LoanAmount:     loanAmount,
			InterestAmount: totalInterest,
			LoanDuration:   loanDuration,
		}

		if err := tx.Create(createdLoan).Error; err != nil {
			return err
		}

		weeklyPrincipal := loanAmount / loanDuration
		weeklyInterest := totalInterest / loanDuration
		weeklyInstallment := weeklyPrincipal + weeklyInterest

		for paymentNumber := 0; paymentNumber < loanDuration; paymentNumber++ {
			var dueDate time.Time

			if paymentNumber == 0 {
				dueDate = firstPaymentDate
			} else {
				dueDate = firstPaymentDate.AddDate(0, 0, 7*paymentNumber)
			}

			if paymentNumber+1 == loanDuration {
				principalDeviation := loanAmount - (weeklyPrincipal * loanDuration)
				interestDeviation := totalInterest - (weeklyInterest * loanDuration)

				weeklyPrincipal += principalDeviation
				weeklyInterest += interestDeviation
				weeklyInstallment = weeklyPrincipal + weeklyInterest
			}

			payment := models.AccountPayment{
				LoanID:          createdLoan.ID,
				DueDate:         dueDate,
				DueAmount:       weeklyInstallment,
				PrincipalAmount: weeklyPrincipal,
				InterestAmount:  weeklyInterest,
			}

			if err := tx.Create(&payment).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdLoan, nil
}

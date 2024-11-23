package main

import (
	"amartha_coding/src/account_payment"
	"amartha_coding/src/loan"
	"amartha_coding/src/models"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

func main() {
	db, err := gorm.Open(sqlite.Open("amartha_coding.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	err = db.AutoMigrate(&models.Loan{}, &models.AccountPayment{})
	if err != nil {
		return
	}

	// Create a new loan
	firstPaymentDate := time.Now()
	initLoan, err := loan_engine.NewLoan(db, 5000000, 50, firstPaymentDate)

	if err != nil {
		log.Fatalf("Failed to create loan and payments: %v\n", err)
	} else {
		log.Printf("Loan created successfully: %+v\n", initLoan)
	}

	loanID := uint(1)
	paymentAmount := 110000
	outstanding, err := account_payment_engine.GetOutstanding(db, loanID)
	fmt.Println("total outstanding: ", outstanding)

	isDelinquent, err := account_payment_engine.IsDelinquent(db, loanID)
	fmt.Println("isDelinquent: ", isDelinquent)

	err = account_payment_engine.MakePayment(db, loanID, paymentAmount)
	if err != nil {
		fmt.Printf("Failed to make payment: %v\n\n", err)
	} else {
		fmt.Printf("Payment of %d made successfully for loan %d\n\n", paymentAmount, loanID)
	}

	outstanding, err = account_payment_engine.GetOutstanding(db, loanID)
	fmt.Println("total outstanding: ", outstanding)
}

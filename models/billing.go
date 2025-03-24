package models

import "sync"

// Loan structure to track loan details
type Billing struct {
	ID                int
	TotalWeeks        int
	LoanAmount        float64
	InterestRate      float64
	TotalAmountDue    float64
	WeeklyInstallment float64
	OutstandingAmount float64
	PaymentHistory    []bool // Track payments (true=paid, false=missed)
	CurrentWeeks      int
}

// Store loans in memory
var loans = make(map[int]*Billing)
var loanMutex sync.Mutex // Protect concurrent access

// Create a new loan
func CreateLoan(id int) *Billing {
	loan := &Billing{
		ID:                id,
		TotalWeeks:        50,
		LoanAmount:        5000000,
		InterestRate:      0.10,
		TotalAmountDue:    5000000 * 1.10,
		WeeklyInstallment: (5000000 * 1.10) / 50,
		OutstandingAmount: 5000000 * 1.10,
		PaymentHistory:    make([]bool, 50), // 50 weeks of payments (default false)
		CurrentWeeks:      0,
	}
	loanMutex.Lock()
	loans[id] = loan
	loanMutex.Unlock()
	return loan
}

// Get an existing loan
func GetLoan(id int) *Billing {
	loanMutex.Lock()
	defer loanMutex.Unlock()
	return loans[id]
}

var (
	loanCounter int
	counterLock sync.Mutex
)

// GenerateLoanID creates a unique loan ID
func GenerateLoanID() int {
	counterLock.Lock()
	defer counterLock.Unlock()
	loanCounter++
	return loanCounter
}

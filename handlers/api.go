package handlers

import (
	"billing-go/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateLoan creates a new loan
// @Summary Create a loan
// @Description Create a new loan with a fixed schedule
// @Tags Loans
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Billing
// @Router /loans/create [post]
func CreateLoan(c *gin.Context) {
	loanID := models.GenerateLoanID() // Generate a new ID
	loan := models.CreateLoan(loanID)
	c.JSON(http.StatusOK, loan)
}

// GetOutstanding returns the outstanding amount
// @Summary Get outstanding amount
// @Description Get remaining balance of a loan
// @Tags Loans
// @Accept  json
// @Produce  json
// @Param id path int true "Loan ID"
// @Success 200 {object} map[string]interface{}
// @Router /loans/{id}/outstanding [get]
func GetOutstanding(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	loan := models.GetLoan(id)
	if loan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"outstanding": loan.OutstandingAmount})
}

// IsDelinquent returns true or false
// @Summary Check IsDelinquent
// @Description Check if a loan is delinquent
// @Tags Loans
// @Accept  json
// @Produce  json
// @Param id path int true "Loan ID"
// @Success 200 {object} map[string]interface{}
// @Router /loans/{id}/delinquent [get]
func IsDelinquent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	loan := models.GetLoan(id)
	if loan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	// Check last two weeks' payments
	unpaidCount := 0
	for i := loan.CurrentWeeks; i >= 0 && unpaidCount < 2; i-- {
		if !loan.PaymentHistory[i] {
			unpaidCount++
		} else {
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"delinquent": unpaidCount >= 2})
}

// MakePayment records a payment
// @Summary Make a payment
// @Description Record a payment
// @Tags Loans
// @Accept  json
// @Produce  json
// @Param id path int true "Loan ID"
// @Success 200 {object} map[string]interface{}
// @Router /loans/{id}/pay [post]
func MakePayment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	loan := models.GetLoan(id)
	if loan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	if loan.OutstandingAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Loan already fully paid"})
		return
	}

	// Find next unpaid week
	for i := range loan.PaymentHistory {
		if !loan.PaymentHistory[i] {
			loan.PaymentHistory[i] = true
			loan.OutstandingAmount -= loan.WeeklyInstallment
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment recorded", "outstanding": loan.OutstandingAmount})
}

// MoveWeeks moves to next week
// @Summary Move to next week
// @Description Move to next week
// @Tags Loans
// @Accept  json
// @Produce  json
// @Param id path int true "Loan ID"
// @Success 200 {object} map[string]interface{}
// @Router /loans/{id}/move-weeks [post]
func MoveWeeks(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	loan := models.GetLoan(id)
	if loan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	if loan.OutstandingAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Loan already fully paid"})
		return
	}

	// Find next unpaid week
	loan.CurrentWeeks += 1

	c.JSON(http.StatusOK, gin.H{"message": "Moves to next week", "currentWeeks": loan.CurrentWeeks})
}

// GetLoan returns the loan
// @Summary Get a loan
// @Description Get a loan
// @Tags Loans
// @Accept  json
// @Produce  json
// @Param id path int true "Loan ID"
// @Success 200 {object} map[string]interface{}
// @Router /loans/{id} [get]
func GetLoan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	loan := models.GetLoan(id)
	if loan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"outstanding": loan})
}

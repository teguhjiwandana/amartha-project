package main

import (
	"billing-go/handlers"

	_ "billing-go/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Loan Billing API
// @version 1.0
// @description This is a loan billing system API.
// @host localhost:8080
// @BasePath /

func main() {
	r := gin.Default() // Buat instance router

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	r.POST("/loans/create", handlers.CreateLoan)
	r.GET("/loans/:id/outstanding", handlers.GetOutstanding)
	r.GET("/loans/:id/delinquent", handlers.IsDelinquent)
	r.POST("/loans/:id/pay", handlers.MakePayment)
	r.POST("/loans/:id/move-weeks", handlers.MoveWeeks)

	r.Run(":8080") // Jalankan server di port 8080
}

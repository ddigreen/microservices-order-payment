package main

import (
	"database/sql"
	"log"

	"payment-service/internal/repository"
	"payment-service/internal/transport/http"
	"payment-service/internal/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "host=localhost port=5432 user=amangeldievdiasbek dbname=payment_db sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	repo := repository.NewSQLPaymentRepository(db)
	uc := usecase.NewPaymentUseCase(repo)
	handler := http.NewPaymentHandler(uc)

	r := gin.Default()

	r.POST("/payments", handler.ProcessPayment)

	r.Run(":8081")
}

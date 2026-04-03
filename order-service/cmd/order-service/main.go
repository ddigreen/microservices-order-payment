package main

import (
	"database/sql"
	"log"

	"order-service/internal/repository"
	"order-service/internal/transport/http"
	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "host=localhost port=5432 user=amangeldievdiasbek dbname=order_db sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Database is unreachable: ", err)
	}

	payClient := repository.NewPaymentClient("http://localhost:8081")
	orderRepo := repository.NewSQLOrderRepository(db)
	orderUC := usecase.NewOrderUseCase(orderRepo, payClient)
	handler := http.NewOrderHandler(orderUC)

	r := gin.Default()

	r.POST("/orders", handler.CreateOrder)
	r.GET("/orders/recent", handler.GetRecent)
	r.GET("/orders/:id", handler.GetOrder)
	r.PATCH("/orders/:id/cancel", handler.CancelOrder)

	log.Println("Order Service starting on :8080...")
	r.Run(":8080")
}

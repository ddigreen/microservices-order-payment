package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"order-service/internal/repository"
	ordergrpc "order-service/internal/transport/grpc"
	"order-service/internal/transport/http"
	"order-service/internal/usecase"

	pb "github.com/ddigreen/payment-generated/payment"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
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

	paymentAddr := os.Getenv("PAYMENT_GRPC_ADDR")
	if paymentAddr == "" {
		paymentAddr = "localhost:50051"
	}

	payClient, err := repository.NewPaymentClient(paymentAddr)
	if err != nil {
		log.Fatalf("Failed to initialize payment client: %v", err)
	}
	defer payClient.Close()

	orderRepo := repository.NewSQLOrderRepository(db)
	orderUC := usecase.NewOrderUseCase(orderRepo, payClient)
	handler := http.NewOrderHandler(orderUC)

	r := gin.Default()

	r.POST("/orders", handler.CreateOrder)
	r.GET("/orders/recent", handler.GetRecent)
	r.GET("/orders/:id", handler.GetOrder)
	r.GET("/payments", handler.GetPayments)
	r.PATCH("/orders/:id/cancel", handler.CancelOrder)

	grpcPort := ":50052"
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port: %v", err)
	}

	grpcServer := grpc.NewServer()
	orderStreamHandler := ordergrpc.NewOrderGrpcServer(orderUC)
	pb.RegisterPaymentServiceServer(grpcServer, orderStreamHandler)

	go func() {
		log.Printf("Order Service gRPC Streaming Server starting on %s...", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to run gRPC server: %v", err)
		}
	}()

	log.Println("Order Service REST HTTP starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run HTTP server: ", err)
	}
}

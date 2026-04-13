package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"payment-service/internal/repository"
	grpchandler "payment-service/internal/transport/grpc"
	"payment-service/internal/usecase"

	pb "github.com/ddigreen/payment-generated/payment"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("PAYMENT_GRPC_PORT")
	if port == "" {
		port = "50051"
	}

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

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	paymentHandler := grpchandler.NewPaymentServer(uc)
	pb.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	log.Printf("Payment gRPC Server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

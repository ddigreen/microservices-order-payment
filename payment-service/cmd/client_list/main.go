package main

import (
	"context"
	"log"

	pb "github.com/ddigreen/payment-generated/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer conn.Close()

	client := pb.NewPaymentServiceClient(conn)

	req := &pb.ListPaymentsRequest{
		MinAmount: 1000,
		MaxAmount: 6000,
	}

	res, err := client.ListPayments(context.Background(), req)
	if err != nil {
		log.Fatalf("Ошибка при запросе: %v", err)
	}

	log.Printf("Найдено платежей: %d", len(res.Payments))
	for i, p := range res.Payments {
		log.Printf("Платеж %d: Статус = %s", i+1, p.Status)
	}
}

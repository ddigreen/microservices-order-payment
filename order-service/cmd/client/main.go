package main

import (
	"context"
	"io"
	"log"
	"os"

	pb "github.com/ddigreen/payment-generated/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Укажи ID заказа! Пример: go run cmd/client/main.go <твой-id-заказа>")
	}
	orderID := os.Args[1]

	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer conn.Close()

	client := pb.NewPaymentServiceClient(conn)

	req := &pb.OrderRequest{OrderId: orderID}
	stream, err := client.SubscribeToOrderUpdates(context.Background(), req)
	if err != nil {
		log.Fatalf("Ошибка подписки: %v", err)
	}

	log.Printf("✅ Подписались на обновления заказа: %s", orderID)
	log.Println("Ожидание изменений статуса...")

	for {
		update, err := stream.Recv()
		if err == io.EOF {
			log.Println("🛑 Стрим закрыт сервером (Достигнут финальный статус).")
			break
		}
		if err != nil {
			log.Fatalf("Ошибка при чтении стрима: %v", err)
		}

		log.Printf(">>> 🔔 НОВЫЙ СТАТУС: %s (Время: %s)", update.Status, update.UpdatedAt.AsTime().Format("15:04:05"))
	}
}

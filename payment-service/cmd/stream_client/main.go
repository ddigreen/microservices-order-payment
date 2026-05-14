package main

import (
	"context"
	"io"
	"log"

	pb "github.com/ddigreen/payment-generated/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer conn.Close()

	client := pb.NewPaymentServiceClient(conn)

	req := &pb.OrderRequest{
		OrderId: "11a27c77-718b-4ac7-b285-2bfbc2917f53",
	}

	stream, err := client.SubscribeToOrderUpdates(context.Background(), req)
	if err != nil {
		log.Fatalf("Ошибка подписки: %v", err)
	}

	log.Printf(" Успешно подписались на заказ: %s", req.OrderId)
	log.Println(" Ожидание обновлений статуса в реальном времени...")

	for {
		update, err := stream.Recv()

		if err == io.EOF {
			log.Println("Сервер закрыл стрим (достигнут финальный статус).")
			break
		}
		if err != nil {
			log.Fatalf("Ошибка при чтении стрима: %v", err)
		}

		log.Printf("[НОВОЕ СОБЫТИЕ] Статус изменился на: %s", update.Status)
	}
}

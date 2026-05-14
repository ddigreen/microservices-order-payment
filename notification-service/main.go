package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentEvent struct {
	OrderID       string `json:"order_id"`
	Amount        int64  `json:"amount"`
	CustomerEmail string `json:"customer_email"`
	Status        string `json:"status"`
}

var (
	processedMessages = make(map[string]bool)
	mutex             sync.Mutex
)

func main() {
	rabbitURL := "amqp://guest:guest@rabbitmq:5672/"
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Failed to open a channel: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare("dlx_exchange", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare DLX: %v", err)
	}

	_, err = ch.QueueDeclare("payment_dlq", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare DLQ: %v", err)
	}

	err = ch.QueueBind("payment_dlq", "", "dlx_exchange", false, nil)
	if err != nil {
		log.Fatalf("Failed to bind DLQ: %v", err)
	}

	args := amqp.Table{
		"x-dead-letter-exchange": "dlx_exchange",
	}

	q, err := ch.QueueDeclare(
		"payment.completed",
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		log.Fatalf("❌ Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("❌ Failed to register a consumer: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for d := range msgs {
			var event PaymentEvent
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Printf("⚠️ Ошибка парсинга JSON: %v", err)
				d.Nack(false, false)
				continue
			}

			if event.Amount == 666 {
				log.Printf("💀 Критическая ошибка! Отправляем заказ %s в мусорку (DLQ)", event.OrderID)
				d.Nack(false, false)
				continue
			}

			mutex.Lock()
			if processedMessages[event.OrderID] {
				log.Printf("⏩ [Idempotency] Сообщение для заказа %s уже обработано. Пропускаем.", event.OrderID)
				mutex.Unlock()
				d.Ack(false)
				continue
			}
			processedMessages[event.OrderID] = true
			mutex.Unlock()

			log.Printf("[Notification] Sent email to %s for Order #%s. Amount: $%d", event.CustomerEmail, event.OrderID, event.Amount)

			err = d.Ack(false)
			if err != nil {
				log.Printf("⚠️ Ошибка ACK: %v", err)
			}
		}
	}()

	log.Println("📥 Notification Service запущен. Ожидание сообщений...")

	<-quit
	log.Println("🛑 Получен сигнал завершения. Выполняем Graceful Shutdown...")
}

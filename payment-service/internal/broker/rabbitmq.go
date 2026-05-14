package broker

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentEvent struct {
	OrderID       string `json:"order_id"`
	Amount        int64  `json:"amount"`
	CustomerEmail string `json:"customer_email"`
	Status        string `json:"status"`
}

type EventPublisher interface {
	PublishPaymentCompleted(ctx context.Context, event PaymentEvent) error
	Close() error
}

type rabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitMQPublisher(amqpURL string) (EventPublisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &rabbitMQPublisher{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (p *rabbitMQPublisher) PublishPaymentCompleted(ctx context.Context, event PaymentEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.channel.PublishWithContext(ctx,
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		})
	if err != nil {
		log.Printf("❌ Ошибка отправки в RabbitMQ: %v", err)
		return err
	}

	log.Printf("✅ Событие отправлено в брокер для заказа: %s", event.OrderID)
	return nil
}

func (p *rabbitMQPublisher) Close() error {
	p.channel.Close()
	return p.conn.Close()
}

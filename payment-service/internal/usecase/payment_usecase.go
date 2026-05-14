package usecase

import (
	"context"
	"errors"
	"log"

	"payment-service/internal/broker"
	"payment-service/internal/domain"

	"github.com/google/uuid"
)

type PaymentUseCase struct {
	repo      domain.PaymentRepository
	publisher broker.EventPublisher
}

func NewPaymentUseCase(r domain.PaymentRepository, pub broker.EventPublisher) *PaymentUseCase {
	return &PaymentUseCase{
		repo:      r,
		publisher: pub,
	}
}

func (u *PaymentUseCase) ProcessPayment(ctx context.Context, p *domain.Payment) (string, string, error) {
	if p.Amount > 100000 {
		p.Status = "Declined"
	} else {
		p.Status = "Authorized"
		p.TransactionID = uuid.New().String()
	}

	if err := u.repo.Create(ctx, p); err != nil {
		return "", "", err
	}

	event := broker.PaymentEvent{
		OrderID:       p.OrderID,
		Amount:        p.Amount,
		CustomerEmail: "user@example.com",
		Status:        p.Status,
	}

	if err := u.publisher.PublishPaymentCompleted(ctx, event); err != nil {
		log.Printf("Warning: failed to publish event to RabbitMQ: %v", err)
	}

	return p.TransactionID, p.Status, nil
}

func (u *PaymentUseCase) ListPayments(ctx context.Context, min, max int64) ([]*domain.Payment, error) {
	if min > 0 && max > 0 && min > max {
		return nil, errors.New("min_amount cannot be greater than max_amount")
	}

	return u.repo.FindByAmountRange(ctx, min, max)
}

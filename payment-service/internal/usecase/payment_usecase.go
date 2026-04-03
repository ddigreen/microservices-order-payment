package usecase

import (
	"context"

	"payment-service/internal/domain"

	"github.com/google/uuid"
)

type PaymentUseCase struct {
	repo domain.PaymentRepository
}

func NewPaymentUseCase(r domain.PaymentRepository) *PaymentUseCase {
	return &PaymentUseCase{repo: r}
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

	return p.TransactionID, p.Status, nil
}

package usecase

import (
	"context"
	"errors"

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

func (u *PaymentUseCase) ListPayments(ctx context.Context, min, max int64) ([]*domain.Payment, error) {
	if min > 0 && max > 0 && min > max {
		return nil, errors.New("min_amount cannot be greater than max_amount")
	}

	// Вызываем метод репозитория, который мы только что создали
	return u.repo.FindByAmountRange(ctx, min, max)
}

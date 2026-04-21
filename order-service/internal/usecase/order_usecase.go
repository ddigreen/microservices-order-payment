package usecase

import (
	"context"
	"errors"

	"order-service/internal/domain"

	pb "github.com/ddigreen/payment-generated/payment"
)

type OrderUseCase struct {
	repo          domain.OrderRepository
	paymentClient domain.PaymentServiceClient
}

func NewOrderUseCase(r domain.OrderRepository, p domain.PaymentServiceClient) *OrderUseCase {
	return &OrderUseCase{repo: r, paymentClient: p}
}

func (u *OrderUseCase) CreateOrder(ctx context.Context, order *domain.Order) error {
	if order.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	order.Status = "Pending"
	if err := u.repo.Create(ctx, order); err != nil {
		return err
	}

	_, status, err := u.paymentClient.AuthorizePayment(ctx, order.ID, order.Amount)
	if err != nil {
		u.repo.UpdateStatus(ctx, order.ID, "Failed")
		return err
	}

	finalStatus := "Failed"
	if status == "Authorized" {
		finalStatus = "Paid"
	}

	return u.repo.UpdateStatus(ctx, order.ID, finalStatus)
}

func (u *OrderUseCase) CancelOrder(ctx context.Context, id string) error {
	order, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status == "Paid" {
		return errors.New("cannot cancel a paid order")
	}

	return u.repo.UpdateStatus(ctx, id, "Cancelled")
}

func (u *OrderUseCase) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	return u.repo.GetByID(ctx, id)
}

func (uc *OrderUseCase) GetRecentOrders(ctx context.Context, limit int) ([]*domain.Order, error) {
	if limit < 1 || limit > 100 {
		return nil, errors.New("limit must be between 1 and 100")
	}

	return uc.repo.GetRecent(ctx, limit)
}

func (u *OrderUseCase) ListPayments(ctx context.Context, min, max int64) (*pb.ListPaymentsResponse, error) {
	return u.paymentClient.ListPayments(ctx, min, max)
}

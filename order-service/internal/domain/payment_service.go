package domain

import (
	"context"

	pb "github.com/ddigreen/payment-generated/payment"
)

type PaymentServiceClient interface {
	AuthorizePayment(ctx context.Context, orderID string, amount int64) (string, string, error)
	ListPayments(ctx context.Context, min, max int64) (*pb.ListPaymentsResponse, error)
}

package domain

import "context"

type PaymentServiceClient interface {
	AuthorizePayment(ctx context.Context, orderID string, amount int64) (string, string, error)
}

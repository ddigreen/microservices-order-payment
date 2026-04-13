package repository

import (
	"context"
	"fmt"
	"time"

	pb "github.com/ddigreen/payment-generated/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type paymentClient struct {
	client pb.PaymentServiceClient
	conn   *grpc.ClientConn
}

func NewPaymentClient(grpcAddress string) (*paymentClient, error) {
	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment grpc server: %w", err)
	}

	client := pb.NewPaymentServiceClient(conn)

	return &paymentClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *paymentClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *paymentClient) AuthorizePayment(ctx context.Context, orderID string, amount int64) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	req := &pb.PaymentRequest{
		OrderId: orderID,
		Amount:  float64(amount),
	}

	resp, err := c.client.ProcessPayment(ctx, req)
	if err != nil {
		return "", "", fmt.Errorf("grpc payment service unavailable: %w", err)
	}

	return "", resp.Status, nil
}

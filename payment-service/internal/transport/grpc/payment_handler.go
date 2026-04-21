package grpc

import (
	"context"
	"log"

	"payment-service/internal/domain"
	"payment-service/internal/usecase"

	pb "github.com/ddigreen/payment-generated/payment"
	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type paymentServer struct {
	pb.UnimplementedPaymentServiceServer
	useCase *usecase.PaymentUseCase
}

func NewPaymentServer(uc *usecase.PaymentUseCase) pb.PaymentServiceServer {
	return &paymentServer{useCase: uc}
}

func (s *paymentServer) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	log.Printf("Received payment request for Order: %s, Amount: %f", req.OrderId, req.Amount)

	paymentInput := &domain.Payment{
		ID:      uuid.New().String(),
		OrderID: req.OrderId,
		Amount:  int64(req.Amount),
	}

	_, paymentStatus, err := s.useCase.ProcessPayment(ctx, paymentInput)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process payment: %v", err)
	}

	return &pb.PaymentResponse{
		Status: paymentStatus,
	}, nil
}

func (s *paymentServer) SubscribeToOrderUpdates(req *pb.OrderRequest, stream pb.PaymentService_SubscribeToOrderUpdatesServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeToOrderUpdates is handled by Order Service")
}

func (s *paymentServer) ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
	payments, err := s.useCase.ListPayments(ctx, req.MinAmount, req.MaxAmount)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	var pbPayments []*pb.PaymentResponse
	for _, p := range payments {
		pbPayments = append(pbPayments, &pb.PaymentResponse{
			Status:  p.Status,
			Id:      p.ID,
			OrderId: p.OrderID,
			Amount:  p.Amount,
		})
	}

	return &pb.ListPaymentsResponse{
		Payments: pbPayments,
	}, nil
}

package grpc

import (
	"context"
	"log"
	"time"

	"order-service/internal/usecase"

	pb "github.com/ddigreen/payment-generated/payment"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type orderGrpcServer struct {
	pb.UnimplementedPaymentServiceServer
	useCase *usecase.OrderUseCase
}

func NewOrderGrpcServer(uc *usecase.OrderUseCase) pb.PaymentServiceServer {
	return &orderGrpcServer{useCase: uc}
}

func (s *orderGrpcServer) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessPayment is handled by Payment Service")
}

func (s *orderGrpcServer) SubscribeToOrderUpdates(req *pb.OrderRequest, stream pb.PaymentService_SubscribeToOrderUpdatesServer) error {
	log.Printf("Client subscribed to updates for Order: %s", req.OrderId)

	lastStatus := ""

	for {

		order, err := s.useCase.GetByID(stream.Context(), req.OrderId)
		if err != nil {
			return status.Errorf(codes.NotFound, "order not found or db error: %v", err)
		}

		if order.Status != lastStatus {
			log.Printf("Pushing new status to stream: %s", order.Status)

			err := stream.Send(&pb.OrderStatusUpdate{
				Status:    order.Status,
				UpdatedAt: timestamppb.Now(),
			})
			if err != nil {
				log.Printf("Client disconnected or stream error: %v", err)
				return err
			}

			lastStatus = order.Status

			if order.Status == "Paid" || order.Status == "Failed" || order.Status == "Cancelled" {
				log.Printf("Final status reached. Closing stream for Order: %s", req.OrderId)
				break
			}
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

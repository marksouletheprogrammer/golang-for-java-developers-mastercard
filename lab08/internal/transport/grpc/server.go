package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"lab08/internal/domain"
	"lab08/internal/repository"
	"lab08/internal/service"
	pb "lab08/proto/orders"
)

// OrderServer implements the gRPC OrderService interface.
// Transport layer - converts between protobuf messages and domain entities.
// Delegates business logic to the service layer (same as HTTP).
type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	service *service.OrderService
}

// NewOrderServer creates a new gRPC server with injected service.
func NewOrderServer(service *service.OrderService) *OrderServer {
	return &OrderServer{
		service: service,
	}
}

// CreateOrder handles gRPC CreateOrder requests.
// TODO: Part 5 - Implement gRPC CreateOrder
func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// TODO: Convert protobuf request to domain.Order using protoToLineItems()
	// TODO: Call s.service.CreateOrder()
	// TODO: Handle errors using mapServiceError()
	// TODO: Convert domain.Order to protobuf using orderToProto()
	// TODO: Return CreateOrderResponse
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// GetOrder handles gRPC GetOrder requests.
// TODO: Part 5 - Implement gRPC GetOrder
func (s *OrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	// TODO: Call s.service.GetOrder() with req.Id
	// TODO: Handle errors using mapServiceError()
	// TODO: Convert order to protobuf and return
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ListOrders handles gRPC ListOrders requests.
// TODO: Part 5 - Implement gRPC ListOrders
func (s *OrderServer) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	// TODO: Call s.service.ListOrders()
	// TODO: Handle errors using mapServiceError()
	// TODO: Convert all orders to protobuf (loop and call orderToProto())
	// TODO: Return ListOrdersResponse with orders array
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// UpdateOrderStatus handles gRPC UpdateOrderStatus requests.
// TODO: Part 5 - Implement gRPC UpdateOrderStatus
func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	// TODO: Convert protobuf status to domain status using protoToStatus()
	// TODO: Call s.service.UpdateOrderStatus()
	// TODO: Handle errors using mapServiceError()
	// TODO: Return empty UpdateOrderStatusResponse on success
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// DeleteOrder handles gRPC DeleteOrder requests.
// TODO: Part 5 - Implement gRPC DeleteOrder
func (s *OrderServer) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	// TODO: Call s.service.DeleteOrder() with req.Id
	// TODO: Handle errors using mapServiceError()
	// TODO: Return empty DeleteOrderResponse on success
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// orderToProto converts domain Order to protobuf Order.
// TODO: Part 5 - Implement conversion helper
func orderToProto(order *domain.Order) *pb.Order {
	// TODO: Create pb.Order and populate fields
	// TODO: Use lineItemsToProto() for items
	// TODO: Use statusToProto() for status
	// TODO: Convert timestamps to Unix (order.CreatedAt.Unix())
	return nil
}

// lineItemsToProto converts domain LineItems to protobuf LineItems.
// TODO: Part 5 - Implement conversion helper
func lineItemsToProto(items []domain.LineItem) []*pb.LineItem {
	// TODO: Create slice of pb.LineItem
	// TODO: Loop through items and convert each one
	// TODO: Note: convert Quantity to int32
	return nil
}

// protoToLineItems converts protobuf LineItems to domain LineItems.
// TODO: Part 5 - Implement conversion helper
func protoToLineItems(items []*pb.LineItem) []domain.LineItem {
	// TODO: Create slice of domain.LineItem
	// TODO: Loop through items and convert each one
	// TODO: Note: convert Quantity from int32 to int
	return nil
}

// statusToProto converts domain OrderStatus to protobuf OrderStatus.
// TODO: Part 5 - Implement conversion helper
func statusToProto(status domain.OrderStatus) pb.OrderStatus {
	// TODO: Use switch statement to map domain status constants to pb.OrderStatus constants
	// TODO: Handle all cases: pending, confirmed, shipped, delivered, cancelled
	return pb.OrderStatus_PENDING
}

// protoToStatus converts protobuf OrderStatus to domain OrderStatus.
// TODO: Part 5 - Implement conversion helper
func protoToStatus(status pb.OrderStatus) domain.OrderStatus {
	// TODO: Use switch statement to map pb.OrderStatus constants to domain status constants
	// TODO: Handle all cases
	return domain.StatusPending
}

// mapServiceError converts service errors to gRPC status codes.
// This is how we communicate errors to gRPC clients.
func mapServiceError(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return status.Error(codes.NotFound, "order not found")
	}
	if errors.Is(err, repository.ErrAlreadyExists) {
		return status.Error(codes.AlreadyExists, "order already exists")
	}
	if errors.Is(err, service.ErrInvalidOrder) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if errors.Is(err, service.ErrInvalidStatusTransition) {
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	return status.Error(codes.Internal, "internal server error")
}

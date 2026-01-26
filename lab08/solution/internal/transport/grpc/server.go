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
func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// Convert protobuf message to domain entity
	order := &domain.Order{
		ID:         req.Id,
		CustomerID: req.CustomerId,
		Items:      protoToLineItems(req.Items),
	}
	
	// Delegate to service layer
	if err := s.service.CreateOrder(ctx, order); err != nil {
		return nil, mapServiceError(err)
	}
	
	// Convert domain entity to protobuf message
	return &pb.CreateOrderResponse{
		Order: orderToProto(order),
	}, nil
}

// GetOrder handles gRPC GetOrder requests.
func (s *OrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := s.service.GetOrder(ctx, req.Id)
	if err != nil {
		return nil, mapServiceError(err)
	}
	
	return &pb.GetOrderResponse{
		Order: orderToProto(order),
	}, nil
}

// ListOrders handles gRPC ListOrders requests.
func (s *OrderServer) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.service.ListOrders(ctx)
	if err != nil {
		return nil, mapServiceError(err)
	}
	
	protoOrders := make([]*pb.Order, len(orders))
	for i, order := range orders {
		protoOrders[i] = orderToProto(order)
	}
	
	return &pb.ListOrdersResponse{
		Orders: protoOrders,
	}, nil
}

// UpdateOrderStatus handles gRPC UpdateOrderStatus requests.
func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	domainStatus := protoToStatus(req.Status)
	
	if err := s.service.UpdateOrderStatus(ctx, req.Id, domainStatus); err != nil {
		return nil, mapServiceError(err)
	}
	
	return &pb.UpdateOrderStatusResponse{}, nil
}

// DeleteOrder handles gRPC DeleteOrder requests.
func (s *OrderServer) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	if err := s.service.DeleteOrder(ctx, req.Id); err != nil {
		return nil, mapServiceError(err)
	}
	
	return &pb.DeleteOrderResponse{}, nil
}

// orderToProto converts domain Order to protobuf Order.
func orderToProto(order *domain.Order) *pb.Order {
	return &pb.Order{
		Id:          order.ID,
		CustomerId:  order.CustomerID,
		Items:       lineItemsToProto(order.Items),
		Status:      statusToProto(order.Status),
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt.Unix(),
		UpdatedAt:   order.UpdatedAt.Unix(),
	}
}

// lineItemsToProto converts domain LineItems to protobuf LineItems.
func lineItemsToProto(items []domain.LineItem) []*pb.LineItem {
	protoItems := make([]*pb.LineItem, len(items))
	for i, item := range items {
		protoItems[i] = &pb.LineItem{
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    int32(item.Quantity),
			UnitPrice:   item.UnitPrice,
		}
	}
	return protoItems
}

// protoToLineItems converts protobuf LineItems to domain LineItems.
func protoToLineItems(items []*pb.LineItem) []domain.LineItem {
	domainItems := make([]domain.LineItem, len(items))
	for i, item := range items {
		domainItems[i] = domain.LineItem{
			ProductID:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    int(item.Quantity),
			UnitPrice:   item.UnitPrice,
		}
	}
	return domainItems
}

// statusToProto converts domain OrderStatus to protobuf OrderStatus.
func statusToProto(status domain.OrderStatus) pb.OrderStatus {
	switch status {
	case domain.StatusPending:
		return pb.OrderStatus_PENDING
	case domain.StatusConfirmed:
		return pb.OrderStatus_CONFIRMED
	case domain.StatusShipped:
		return pb.OrderStatus_SHIPPED
	case domain.StatusDelivered:
		return pb.OrderStatus_DELIVERED
	case domain.StatusCancelled:
		return pb.OrderStatus_CANCELLED
	default:
		return pb.OrderStatus_PENDING
	}
}

// protoToStatus converts protobuf OrderStatus to domain OrderStatus.
func protoToStatus(status pb.OrderStatus) domain.OrderStatus {
	switch status {
	case pb.OrderStatus_PENDING:
		return domain.StatusPending
	case pb.OrderStatus_CONFIRMED:
		return domain.StatusConfirmed
	case pb.OrderStatus_SHIPPED:
		return domain.StatusShipped
	case pb.OrderStatus_DELIVERED:
		return domain.StatusDelivered
	case pb.OrderStatus_CANCELLED:
		return domain.StatusCancelled
	default:
		return domain.StatusPending
	}
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

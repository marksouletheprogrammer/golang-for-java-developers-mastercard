package http

import (
	"encoding/json"
	"net/http"

	"lab08/internal/domain"
	"lab08/internal/service"
)

// OrderHandler handles HTTP requests for order operations.
// Transport layer - focuses on HTTP concerns (request/response handling).
// Delegates business logic to the service layer.
type OrderHandler struct {
	service *service.OrderService
}

// NewOrderHandler creates a new HTTP handler with injected service.
func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

// CreateOrderRequest represents the JSON structure for creating orders.
type CreateOrderRequest struct {
	ID         string            `json:"id"`
	CustomerID string            `json:"customer_id"`
	Items      []domain.LineItem `json:"items"`
}

// UpdateStatusRequest represents the JSON structure for status updates.
type UpdateStatusRequest struct {
	Status domain.OrderStatus `json:"status"`
}

// CreateOrder handles POST /orders - creates a new order.
// TODO: Part 3 - Implement HTTP CreateOrder handler
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// TODO: Check method is POST
	// TODO: Decode JSON body into CreateOrderRequest
	// TODO: Create domain.Order from request
	// TODO: Call h.service.CreateOrder()
	// TODO: Handle errors appropriately:
	//   - service.ErrInvalidOrder -> 400 Bad Request
	//   - repository.ErrAlreadyExists -> 409 Conflict
	//   - Other errors -> 500 Internal Server Error
	// TODO: Respond with created order and 201 status
}

// GetOrder handles GET /orders/{id} - retrieves an order by ID.
// TODO: Part 3 - Implement HTTP GetOrder handler
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// TODO: Check method is GET
	// TODO: Extract ID from URL path
	// TODO: Call h.service.GetOrder()
	// TODO: Handle errors:
	//   - repository.ErrNotFound -> 404 Not Found
	//   - Other errors -> 500 Internal Server Error
	// TODO: Respond with order and 200 status
}

// ListOrders handles GET /orders - retrieves all orders.
// TODO: Part 3 - Implement HTTP ListOrders handler
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	// TODO: Check method is GET
	// TODO: Call h.service.ListOrders()
	// TODO: Handle errors (500 for any error)
	// TODO: Respond with orders array and 200 status
}

// UpdateOrderStatus handles PATCH /orders/{id}/status - updates order status.
// TODO: Part 3 - Implement HTTP UpdateOrderStatus handler
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	// TODO: Check method is PATCH
	// TODO: Extract ID from URL path (handle /orders/{id}/status format)
	// TODO: Decode JSON body into UpdateStatusRequest
	// TODO: Call h.service.UpdateOrderStatus()
	// TODO: Handle errors:
	//   - repository.ErrNotFound -> 404 Not Found
	//   - service.ErrInvalidStatusTransition -> 400 Bad Request
	//   - Other errors -> 500 Internal Server Error
	// TODO: Respond with 204 No Content on success
}

// DeleteOrder handles DELETE /orders/{id} - deletes an order.
// TODO: Part 3 - Implement HTTP DeleteOrder handler
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	// TODO: Check method is DELETE
	// TODO: Extract ID from URL path
	// TODO: Call h.service.DeleteOrder()
	// TODO: Handle errors:
	//   - repository.ErrNotFound -> 404 Not Found
	//   - Other errors -> 500 Internal Server Error
	// TODO: Respond with 204 No Content on success
}

// respondJSON writes a JSON response with the given status code.
func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// respondError writes an error response in JSON format.
func respondError(w http.ResponseWriter, message string, statusCode int) {
	respondJSON(w, ErrorResponse{Error: message}, statusCode)
}

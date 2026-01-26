package http

import (
	"encoding/json"
	"errors"
	"net/http"

	commonhttp "golang-for-java-developers-training/common/http"
	"lab08/internal/domain"
	"lab08/internal/repository"
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
	ID         string               `json:"id"`
	CustomerID string               `json:"customer_id"`
	Items      []domain.LineItem    `json:"items"`
}

// UpdateStatusRequest represents the JSON structure for status updates.
type UpdateStatusRequest struct {
	Status domain.OrderStatus `json:"status"`
}

// CreateOrder handles POST /orders - creates a new order.
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	order := &domain.Order{
		ID:         req.ID,
		CustomerID: req.CustomerID,
		Items:      req.Items,
	}
	
	if err := h.service.CreateOrder(r.Context(), order); err != nil {
		if errors.Is(err, service.ErrInvalidOrder) {
			respondError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, repository.ErrAlreadyExists) {
			respondError(w, "Order already exists", http.StatusConflict)
			return
		}
		respondError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, order, http.StatusCreated)
}

// GetOrder handles GET /orders/{id} - retrieves an order by ID.
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Extract ID from URL path (simple path parsing)
	id := r.URL.Path[len("/orders/"):]
	if id == "" {
		respondError(w, "Order ID required", http.StatusBadRequest)
		return
	}
	
	order, err := h.service.GetOrder(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, "Order not found", http.StatusNotFound)
			return
		}
		respondError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, order, http.StatusOK)
}

// ListOrders handles GET /orders - retrieves all orders.
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	orders, err := h.service.ListOrders(r.Context())
	if err != nil {
		respondError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, orders, http.StatusOK)
}

// UpdateOrderStatus handles PATCH /orders/{id}/status - updates order status.
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Extract ID from URL path
	id := r.URL.Path[len("/orders/"):]
	id = id[:len(id)-len("/status")]
	if id == "" {
		respondError(w, "Order ID required", http.StatusBadRequest)
		return
	}
	
	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if err := h.service.UpdateOrderStatus(r.Context(), id, req.Status); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, "Order not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrInvalidStatusTransition) {
			respondError(w, err.Error(), http.StatusBadRequest)
			return
		}
		respondError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// DeleteOrder handles DELETE /orders/{id} - deletes an order.
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Extract ID from URL path
	id := r.URL.Path[len("/orders/"):]
	if id == "" {
		respondError(w, "Order ID required", http.StatusBadRequest)
		return
	}
	
	if err := h.service.DeleteOrder(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, "Order not found", http.StatusNotFound)
			return
		}
		respondError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// respondJSON writes a JSON response with the given status code.
func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// respondError writes an error response in JSON format.
func respondError(w http.ResponseWriter, message string, statusCode int) {
	respondJSON(w, commonhttp.ErrorResponse{Error: message}, statusCode)
}

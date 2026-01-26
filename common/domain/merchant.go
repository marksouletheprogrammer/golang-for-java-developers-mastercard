package domain

import (
	"errors"
	"sync"
)

// Merchant represents a merchant in the system.
type Merchant struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Country  string `json:"country"`
	Status   string `json:"status"`
}

// MerchantStore provides thread-safe in-memory storage for merchants.
// Uses a sync.RWMutex to allow concurrent reads while ensuring exclusive writes.
type MerchantStore struct {
	mu        sync.RWMutex
	merchants map[string]*Merchant
}

// NewMerchantStore creates and initializes a new merchant store with sample data.
func NewMerchantStore() *MerchantStore {
	store := &MerchantStore{
		merchants: make(map[string]*Merchant),
	}

	// Initialize with sample merchants
	store.merchants["MERCH-001"] = &Merchant{
		ID:       "MERCH-001",
		Name:     "Tech Solutions Inc",
		Category: "Technology",
		Country:  "USA",
		Status:   "active",
	}
	store.merchants["MERCH-002"] = &Merchant{
		ID:       "MERCH-002",
		Name:     "Global Retail Co",
		Category: "Retail",
		Country:  "UK",
		Status:   "active",
	}
	store.merchants["MERCH-003"] = &Merchant{
		ID:       "MERCH-003",
		Name:     "Food Services Ltd",
		Category: "Food & Beverage",
		Country:  "Canada",
		Status:   "inactive",
	}

	return store
}

// GetAll returns all merchants.
// Uses RLock for concurrent read access - multiple goroutines can read simultaneously.
func (s *MerchantStore) GetAll() []*Merchant {
	s.mu.RLock()
	defer s.mu.RUnlock()

	merchants := make([]*Merchant, 0, len(s.merchants))
	for _, merchant := range s.merchants {
		merchants = append(merchants, merchant)
	}
	return merchants
}

// GetByID returns a merchant by ID or an error if not found.
// Uses RLock for concurrent read access.
func (s *MerchantStore) GetByID(id string) (*Merchant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	merchant, exists := s.merchants[id]
	if !exists {
		return nil, errors.New("merchant not found")
	}
	return merchant, nil
}

// Create adds a new merchant to the store.
// Uses Lock for exclusive write access - no other reads or writes can happen.
// Returns an error if a merchant with the same ID already exists.
func (s *MerchantStore) Create(merchant *Merchant) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if merchant already exists
	if _, exists := s.merchants[merchant.ID]; exists {
		return errors.New("merchant with this ID already exists")
	}

	s.merchants[merchant.ID] = merchant
	return nil
}

// ValidateMerchant checks if a merchant has all required fields.
func ValidateMerchant(m *Merchant) error {
	if m.ID == "" {
		return errors.New("id is required")
	}
	if m.Name == "" {
		return errors.New("name is required")
	}
	if m.Category == "" {
		return errors.New("category is required")
	}
	if m.Country == "" {
		return errors.New("country is required")
	}
	if m.Status == "" {
		return errors.New("status is required")
	}

	// Validate status values
	if m.Status != "active" && m.Status != "inactive" {
		return errors.New("status must be 'active' or 'inactive'")
	}

	return nil
}

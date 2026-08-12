package memory

import (
	"sync"

	"github.com/LICObr/payrest/internal/domain"
)

// ChargeRepository is a thread-safe in-memory implementation of
// domain.ChargeRepository.
type ChargeRepository struct {
	mu      sync.RWMutex
	charges map[string]*domain.Charge
}

// --------------------------------------------------
// NewChargeRepository
// --------------------------------------------------

func NewChargeRepository() *ChargeRepository {
	return &ChargeRepository{
		charges: make(map[string]*domain.Charge),
	}
}

// --------------------------------------------------
// Save
// --------------------------------------------------

func (r *ChargeRepository) Save(charge *domain.Charge) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.charges[charge.ID] = charge
	return nil
}

// --------------------------------------------------
// FindByID
// --------------------------------------------------

func (r *ChargeRepository) FindByID(id string) (*domain.Charge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	charge, ok := r.charges[id]
	if !ok {
		return nil, domain.ErrChargeNotFound
	}
	return charge, nil
}

// --------------------------------------------------
// FindByAccountID
// --------------------------------------------------

func (r *ChargeRepository) FindByAccountID(accountID string) ([]*domain.Charge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Charge
	for _, charge := range r.charges {
		if charge.AccountID == accountID {
			result = append(result, charge)
		}
	}
	return result, nil
}

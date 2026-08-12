package memory

import (
	"sync"

	"github.com/LICObr/payrest/internal/domain"
)

// AccountRepository is a thread-safe in-memory implementation of
// domain.AccountRepository. It is meant for demos and tests; swapping it
// for a real database only requires implementing the same interface.
type AccountRepository struct {
	mu       sync.RWMutex
	accounts map[string]*domain.Account
}

// --------------------------------------------------
// NewAccountRepository
// --------------------------------------------------

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{
		accounts: make(map[string]*domain.Account),
	}
}

// --------------------------------------------------
// Save
// --------------------------------------------------

func (r *AccountRepository) Save(account *domain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.accounts[account.ID] = account
	return nil
}

// --------------------------------------------------
// FindByID
// --------------------------------------------------

func (r *AccountRepository) FindByID(id string) (*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, ok := r.accounts[id]
	if !ok {
		return nil, domain.ErrAccountNotFound
	}
	return account, nil
}

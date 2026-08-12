package service

import "github.com/LICObr/payrest/internal/domain"

// AccountService handles account use cases. It depends only on the
// AccountRepository interface, never on a concrete storage engine.
type AccountService struct {
	accounts domain.AccountRepository
}

// --------------------------------------------------
// NewAccountService
// --------------------------------------------------

func NewAccountService(accounts domain.AccountRepository) *AccountService {
	return &AccountService{accounts: accounts}
}

// --------------------------------------------------
// CreateAccount
// --------------------------------------------------

func (s *AccountService) CreateAccount(name string, initialBalanceCents int64) (*domain.Account, error) {
	if initialBalanceCents < 0 {
		return nil, domain.ErrInvalidAmount
	}

	account := domain.NewAccount(generateID(), name, initialBalanceCents)
	if err := s.accounts.Save(account); err != nil {
		return nil, err
	}
	return account, nil
}

// --------------------------------------------------
// GetAccount
// --------------------------------------------------

func (s *AccountService) GetAccount(accountID string) (*domain.Account, error) {
	return s.accounts.FindByID(accountID)
}

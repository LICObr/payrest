package service

import "github.com/LICObr/payrest/internal/domain"

// ChargeService handles the payment gateway use cases: creating a charge
// and moving it through its state machine (pending -> approved|declined -> refunded).
type ChargeService struct {
	charges  domain.ChargeRepository
	accounts domain.AccountRepository
}

// --------------------------------------------------
// NewChargeService
// --------------------------------------------------

func NewChargeService(charges domain.ChargeRepository, accounts domain.AccountRepository) *ChargeService {
	return &ChargeService{charges: charges, accounts: accounts}
}

// --------------------------------------------------
// CreateCharge
// --------------------------------------------------

func (s *ChargeService) CreateCharge(accountID string, amountCents int64) (*domain.Charge, error) {
	if amountCents <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	if _, err := s.accounts.FindByID(accountID); err != nil {
		return nil, err
	}

	charge := domain.NewCharge(generateID(), accountID, amountCents)
	if err := s.charges.Save(charge); err != nil {
		return nil, err
	}
	return charge, nil
}

// --------------------------------------------------
// GetCharge
// --------------------------------------------------

func (s *ChargeService) GetCharge(chargeID string) (*domain.Charge, error) {
	return s.charges.FindByID(chargeID)
}

package service

import (
	"sync"

	"github.com/LICObr/payrest/internal/domain"
)

// ChargeService handles the payment gateway use cases: creating a charge
// and moving it through its state machine (pending -> approved|declined -> refunded).
//
// State transitions (Approve/Decline/Refund) read a charge, decide whether
// the transition is valid, then write the result back in separate steps.
// Without a lock, two concurrent requests for the same charge could both
// pass the validity check before either writes — e.g. double-debiting an
// account on a double approve. mu serializes those transitions.
type ChargeService struct {
	charges  domain.ChargeRepository
	accounts domain.AccountRepository
	mu       sync.Mutex
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
// Approve
// --------------------------------------------------

func (s *ChargeService) Approve(chargeID string) (*domain.Charge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	charge, err := s.charges.FindByID(chargeID)
	if err != nil {
		return nil, err
	}
	if !charge.CanApprove() {
		return nil, domain.ErrInvalidChargeState
	}

	account, err := s.accounts.FindByID(charge.AccountID)
	if err != nil {
		return nil, err
	}
	if err := account.Debit(charge.AmountCents); err != nil {
		return nil, err
	}

	charge.MarkApproved()
	if err := s.accounts.Save(account); err != nil {
		return nil, err
	}
	if err := s.charges.Save(charge); err != nil {
		return nil, err
	}
	return charge, nil
}

// --------------------------------------------------
// Decline
// --------------------------------------------------

func (s *ChargeService) Decline(chargeID string) (*domain.Charge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	charge, err := s.charges.FindByID(chargeID)
	if err != nil {
		return nil, err
	}
	if !charge.CanDecline() {
		return nil, domain.ErrInvalidChargeState
	}

	charge.MarkDeclined()
	if err := s.charges.Save(charge); err != nil {
		return nil, err
	}
	return charge, nil
}

// --------------------------------------------------
// Refund
// --------------------------------------------------

func (s *ChargeService) Refund(chargeID string) (*domain.Charge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	charge, err := s.charges.FindByID(chargeID)
	if err != nil {
		return nil, err
	}
	if !charge.CanRefund() {
		return nil, domain.ErrInvalidChargeState
	}

	account, err := s.accounts.FindByID(charge.AccountID)
	if err != nil {
		return nil, err
	}
	account.Credit(charge.AmountCents)

	charge.MarkRefunded()
	if err := s.accounts.Save(account); err != nil {
		return nil, err
	}
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

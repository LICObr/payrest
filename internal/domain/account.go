package domain

import "time"

// Account represents a holder of funds inside the gateway.
// Balance is stored in cents to avoid floating point rounding errors.
type Account struct {
	ID           string
	Name         string
	BalanceCents int64
	CreatedAt    time.Time
}

// --------------------------------------------------
// NewAccount
// --------------------------------------------------

func NewAccount(id, name string, initialBalanceCents int64) *Account {
	return &Account{
		ID:           id,
		Name:         name,
		BalanceCents: initialBalanceCents,
		CreatedAt:    time.Now(),
	}
}

// --------------------------------------------------
// Credit
// --------------------------------------------------

func (a *Account) Credit(amountCents int64) {
	a.BalanceCents += amountCents
}

// --------------------------------------------------
// Debit
// --------------------------------------------------

func (a *Account) Debit(amountCents int64) error {
	if amountCents > a.BalanceCents {
		return ErrInsufficientFunds
	}
	a.BalanceCents -= amountCents
	return nil
}

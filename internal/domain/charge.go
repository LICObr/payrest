package domain

import "time"

type ChargeStatus string

const (
	ChargeStatusPending  ChargeStatus = "pending"
	ChargeStatusApproved ChargeStatus = "approved"
	ChargeStatusDeclined ChargeStatus = "declined"
	ChargeStatusRefunded ChargeStatus = "refunded"
)

// Charge represents a single payment attempt against an account.
// It moves through a small state machine: pending -> approved|declined -> refunded.
type Charge struct {
	ID          string
	AccountID   string
	AmountCents int64
	Status      ChargeStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// --------------------------------------------------
// NewCharge
// --------------------------------------------------

func NewCharge(id, accountID string, amountCents int64) *Charge {
	now := time.Now()
	return &Charge{
		ID:          id,
		AccountID:   accountID,
		AmountCents: amountCents,
		Status:      ChargeStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// --------------------------------------------------
// CanApprove
// --------------------------------------------------

func (c *Charge) CanApprove() bool {
	return c.Status == ChargeStatusPending
}

// --------------------------------------------------
// CanDecline
// --------------------------------------------------

func (c *Charge) CanDecline() bool {
	return c.Status == ChargeStatusPending
}

// --------------------------------------------------
// CanRefund
// --------------------------------------------------

func (c *Charge) CanRefund() bool {
	return c.Status == ChargeStatusApproved
}

// --------------------------------------------------
// MarkApproved
// --------------------------------------------------

func (c *Charge) MarkApproved() {
	c.Status = ChargeStatusApproved
	c.UpdatedAt = time.Now()
}

// --------------------------------------------------
// MarkDeclined
// --------------------------------------------------

func (c *Charge) MarkDeclined() {
	c.Status = ChargeStatusDeclined
	c.UpdatedAt = time.Now()
}

// --------------------------------------------------
// MarkRefunded
// --------------------------------------------------

func (c *Charge) MarkRefunded() {
	c.Status = ChargeStatusRefunded
	c.UpdatedAt = time.Now()
}

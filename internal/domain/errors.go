package domain

import "errors"

var (
	ErrInsufficientFunds  = errors.New("account does not have enough balance")
	ErrInvalidAmount      = errors.New("amount must be greater than zero")
	ErrAccountNotFound    = errors.New("account not found")
	ErrChargeNotFound     = errors.New("charge not found")
	ErrInvalidChargeState = errors.New("charge cannot transition from its current state")
)

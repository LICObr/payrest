package handler

import "github.com/LICObr/payrest/internal/domain"

type CreateAccountRequest struct {
	Name                string `json:"name"`
	InitialBalanceCents int64  `json:"initial_balance_cents"`
}

type AccountResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	BalanceCents int64  `json:"balance_cents"`
}

type CreateChargeRequest struct {
	AccountID   string `json:"account_id"`
	AmountCents int64  `json:"amount_cents"`
}

type ChargeResponse struct {
	ID          string `json:"id"`
	AccountID   string `json:"account_id"`
	AmountCents int64  `json:"amount_cents"`
	Status      string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// --------------------------------------------------
// toAccountResponse
// --------------------------------------------------

func toAccountResponse(account *domain.Account) AccountResponse {
	return AccountResponse{
		ID:           account.ID,
		Name:         account.Name,
		BalanceCents: account.BalanceCents,
	}
}

// --------------------------------------------------
// toChargeResponse
// --------------------------------------------------

func toChargeResponse(charge *domain.Charge) ChargeResponse {
	return ChargeResponse{
		ID:          charge.ID,
		AccountID:   charge.AccountID,
		AmountCents: charge.AmountCents,
		Status:      string(charge.Status),
	}
}

// --------------------------------------------------
// toChargeResponses
// --------------------------------------------------

func toChargeResponses(charges []*domain.Charge) []ChargeResponse {
	responses := make([]ChargeResponse, len(charges))
	for i, charge := range charges {
		responses[i] = toChargeResponse(charge)
	}
	return responses
}

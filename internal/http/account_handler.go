package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LICObr/payrest/internal/domain"
	"github.com/LICObr/payrest/internal/service"
)

type AccountHandler struct {
	accounts *service.AccountService
}

// --------------------------------------------------
// NewAccountHandler
// --------------------------------------------------

func NewAccountHandler(accounts *service.AccountService) *AccountHandler {
	return &AccountHandler{accounts: accounts}
}

// --------------------------------------------------
// Create
// --------------------------------------------------

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	account, err := h.accounts.CreateAccount(req.Name, req.InitialBalanceCents)
	if err != nil {
		writeError(w, statusForError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, toAccountResponse(account))
}

// --------------------------------------------------
// Get
// --------------------------------------------------

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	account, err := h.accounts.GetAccount(id)
	if err != nil {
		writeError(w, statusForError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toAccountResponse(account))
}

// --------------------------------------------------
// writeJSON
// --------------------------------------------------

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// --------------------------------------------------
// writeError
// --------------------------------------------------

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

// --------------------------------------------------
// statusForError
// --------------------------------------------------

func statusForError(err error) int {
	switch {
	case errors.Is(err, domain.ErrAccountNotFound), errors.Is(err, domain.ErrChargeNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrInsufficientFunds), errors.Is(err, domain.ErrInvalidChargeState):
		return http.StatusConflict
	case errors.Is(err, domain.ErrInvalidAmount):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

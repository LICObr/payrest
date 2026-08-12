package handler

import (
	"encoding/json"
	"net/http"

	"github.com/LICObr/payrest/internal/service"
)

type ChargeHandler struct {
	charges *service.ChargeService
}

// --------------------------------------------------
// NewChargeHandler
// --------------------------------------------------

func NewChargeHandler(charges *service.ChargeService) *ChargeHandler {
	return &ChargeHandler{charges: charges}
}

// --------------------------------------------------
// Create
// --------------------------------------------------

func (h *ChargeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateChargeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	charge, err := h.charges.CreateCharge(req.AccountID, req.AmountCents)
	if err != nil {
		writeError(w, statusForError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, toChargeResponse(charge))
}

// --------------------------------------------------
// Get
// --------------------------------------------------

func (h *ChargeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	charge, err := h.charges.GetCharge(id)
	if err != nil {
		writeError(w, statusForError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toChargeResponse(charge))
}

// --------------------------------------------------
// Approve
// --------------------------------------------------

func (h *ChargeHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	charge, err := h.charges.Approve(id)
	if err != nil {
		writeError(w, statusForError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toChargeResponse(charge))
}

// --------------------------------------------------
// Decline
// --------------------------------------------------

func (h *ChargeHandler) Decline(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	charge, err := h.charges.Decline(id)
	if err != nil {
		writeError(w, statusForError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toChargeResponse(charge))
}

// --------------------------------------------------
// Refund
// --------------------------------------------------

func (h *ChargeHandler) Refund(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	charge, err := h.charges.Refund(id)
	if err != nil {
		writeError(w, statusForError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toChargeResponse(charge))
}

// --------------------------------------------------
// ListByAccount
// --------------------------------------------------

func (h *ChargeHandler) ListByAccount(w http.ResponseWriter, r *http.Request) {
	accountID := r.PathValue("id")

	charges, err := h.charges.ListByAccount(accountID)
	if err != nil {
		writeError(w, statusForError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toChargeResponses(charges))
}

package handler

import "net/http"

// --------------------------------------------------
// NewRouter
// --------------------------------------------------

func NewRouter(accounts *AccountHandler, charges *ChargeHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handleHealth)

	mux.HandleFunc("POST /accounts", accounts.Create)
	mux.HandleFunc("GET /accounts/{id}", accounts.Get)
	mux.HandleFunc("GET /accounts/{id}/charges", charges.ListByAccount)

	mux.HandleFunc("POST /charges", charges.Create)
	mux.HandleFunc("GET /charges/{id}", charges.Get)
	mux.HandleFunc("POST /charges/{id}/approve", charges.Approve)
	mux.HandleFunc("POST /charges/{id}/decline", charges.Decline)
	mux.HandleFunc("POST /charges/{id}/refund", charges.Refund)

	return mux
}

// --------------------------------------------------
// handleHealth
// --------------------------------------------------

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

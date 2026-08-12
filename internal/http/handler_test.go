package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	payresthttp "github.com/LICObr/payrest/internal/http"
	"github.com/LICObr/payrest/internal/repository/memory"
	"github.com/LICObr/payrest/internal/service"
)

// --------------------------------------------------
// newTestRouter
// --------------------------------------------------

func newTestRouter() http.Handler {
	accountRepo := memory.NewAccountRepository()
	chargeRepo := memory.NewChargeRepository()

	accountService := service.NewAccountService(accountRepo)
	chargeService := service.NewChargeService(chargeRepo, accountRepo)

	accountHandler := payresthttp.NewAccountHandler(accountService)
	chargeHandler := payresthttp.NewChargeHandler(chargeService)

	return payresthttp.NewRouter(accountHandler, chargeHandler)
}

// --------------------------------------------------
// postJSON
// --------------------------------------------------

func postJSON(router http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

// --------------------------------------------------
// TestCreateAccount_ReturnsCreated
// --------------------------------------------------

func TestCreateAccount_ReturnsCreated(t *testing.T) {
	router := newTestRouter()

	rec := postJSON(router, http.MethodPost, "/accounts", payresthttp.CreateAccountRequest{
		Name:                "Danilo",
		InitialBalanceCents: 10_000,
	})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var account payresthttp.AccountResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &account)
	if account.BalanceCents != 10_000 {
		t.Errorf("expected balance 10000, got %d", account.BalanceCents)
	}
}

// --------------------------------------------------
// TestGetAccount_NotFound_Returns404
// --------------------------------------------------

func TestGetAccount_NotFound_Returns404(t *testing.T) {
	router := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/accounts/does-not-exist", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

// --------------------------------------------------
// TestChargeLifecycle_CreateApproveRefund
// --------------------------------------------------

func TestChargeLifecycle_CreateApproveRefund(t *testing.T) {
	router := newTestRouter()

	accountRec := postJSON(router, http.MethodPost, "/accounts", payresthttp.CreateAccountRequest{
		Name:                "Danilo",
		InitialBalanceCents: 10_000,
	})
	var account payresthttp.AccountResponse
	_ = json.Unmarshal(accountRec.Body.Bytes(), &account)

	chargeRec := postJSON(router, http.MethodPost, "/charges", payresthttp.CreateChargeRequest{
		AccountID:   account.ID,
		AmountCents: 3_000,
	})
	if chargeRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 creating charge, got %d: %s", chargeRec.Code, chargeRec.Body.String())
	}
	var charge payresthttp.ChargeResponse
	_ = json.Unmarshal(chargeRec.Body.Bytes(), &charge)

	approveRec := postJSON(router, http.MethodPost, "/charges/"+charge.ID+"/approve", nil)
	if approveRec.Code != http.StatusOK {
		t.Fatalf("expected 200 approving charge, got %d: %s", approveRec.Code, approveRec.Body.String())
	}

	refundRec := postJSON(router, http.MethodPost, "/charges/"+charge.ID+"/refund", nil)
	if refundRec.Code != http.StatusOK {
		t.Fatalf("expected 200 refunding charge, got %d: %s", refundRec.Code, refundRec.Body.String())
	}

	var refunded payresthttp.ChargeResponse
	_ = json.Unmarshal(refundRec.Body.Bytes(), &refunded)
	if refunded.Status != "refunded" {
		t.Errorf("expected status refunded, got %s", refunded.Status)
	}
}

// --------------------------------------------------
// TestListChargesByAccount_ReturnsOnlyThatAccountsCharges
// --------------------------------------------------

func TestListChargesByAccount_ReturnsOnlyThatAccountsCharges(t *testing.T) {
	router := newTestRouter()

	accountARec := postJSON(router, http.MethodPost, "/accounts", payresthttp.CreateAccountRequest{Name: "A", InitialBalanceCents: 10_000})
	var accountA payresthttp.AccountResponse
	_ = json.Unmarshal(accountARec.Body.Bytes(), &accountA)

	accountBRec := postJSON(router, http.MethodPost, "/accounts", payresthttp.CreateAccountRequest{Name: "B", InitialBalanceCents: 10_000})
	var accountB payresthttp.AccountResponse
	_ = json.Unmarshal(accountBRec.Body.Bytes(), &accountB)

	postJSON(router, http.MethodPost, "/charges", payresthttp.CreateChargeRequest{AccountID: accountA.ID, AmountCents: 1_000})
	postJSON(router, http.MethodPost, "/charges", payresthttp.CreateChargeRequest{AccountID: accountA.ID, AmountCents: 2_000})
	postJSON(router, http.MethodPost, "/charges", payresthttp.CreateChargeRequest{AccountID: accountB.ID, AmountCents: 5_000})

	req := httptest.NewRequest(http.MethodGet, "/accounts/"+accountA.ID+"/charges", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var charges []payresthttp.ChargeResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &charges)
	if len(charges) != 2 {
		t.Errorf("expected 2 charges for account A, got %d", len(charges))
	}
}

package service_test

import (
	"testing"

	"github.com/LICObr/payrest/internal/domain"
	"github.com/LICObr/payrest/internal/repository/memory"
	"github.com/LICObr/payrest/internal/service"
)

// --------------------------------------------------
// newGatewayForTest
// --------------------------------------------------

func newGatewayForTest() (*service.AccountService, *service.ChargeService) {
	accountRepo := memory.NewAccountRepository()
	chargeRepo := memory.NewChargeRepository()

	accounts := service.NewAccountService(accountRepo)
	charges := service.NewChargeService(chargeRepo, accountRepo)
	return accounts, charges
}

// --------------------------------------------------
// TestCreateCharge_AccountNotFound
// --------------------------------------------------

func TestCreateCharge_AccountNotFound(t *testing.T) {
	_, charges := newGatewayForTest()

	_, err := charges.CreateCharge("does-not-exist", 1_000)
	if err != domain.ErrAccountNotFound {
		t.Errorf("expected ErrAccountNotFound, got %v", err)
	}
}

// --------------------------------------------------
// TestCreateCharge_InvalidAmount
// --------------------------------------------------

func TestCreateCharge_InvalidAmount(t *testing.T) {
	accounts, charges := newGatewayForTest()
	account, _ := accounts.CreateAccount("Danilo", 10_000)

	_, err := charges.CreateCharge(account.ID, 0)
	if err != domain.ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

// --------------------------------------------------
// TestApprove_DebitsAccountBalance
// --------------------------------------------------

func TestApprove_DebitsAccountBalance(t *testing.T) {
	accounts, charges := newGatewayForTest()
	account, _ := accounts.CreateAccount("Danilo", 10_000)
	charge, _ := charges.CreateCharge(account.ID, 3_000)

	approved, err := charges.Approve(charge.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if approved.Status != domain.ChargeStatusApproved {
		t.Errorf("expected status approved, got %s", approved.Status)
	}

	updatedAccount, _ := accounts.GetAccount(account.ID)
	if updatedAccount.BalanceCents != 7_000 {
		t.Errorf("expected balance 7000 after approval, got %d", updatedAccount.BalanceCents)
	}
}

// --------------------------------------------------
// TestApprove_InsufficientFunds
// --------------------------------------------------

func TestApprove_InsufficientFunds(t *testing.T) {
	accounts, charges := newGatewayForTest()
	account, _ := accounts.CreateAccount("Danilo", 1_000)
	charge, _ := charges.CreateCharge(account.ID, 5_000)

	_, err := charges.Approve(charge.ID)
	if err != domain.ErrInsufficientFunds {
		t.Errorf("expected ErrInsufficientFunds, got %v", err)
	}
}

// --------------------------------------------------
// TestApprove_AlreadyApprovedCannotApproveAgain
// --------------------------------------------------

func TestApprove_AlreadyApprovedCannotApproveAgain(t *testing.T) {
	accounts, charges := newGatewayForTest()
	account, _ := accounts.CreateAccount("Danilo", 10_000)
	charge, _ := charges.CreateCharge(account.ID, 3_000)
	_, _ = charges.Approve(charge.ID)

	_, err := charges.Approve(charge.ID)
	if err != domain.ErrInvalidChargeState {
		t.Errorf("expected ErrInvalidChargeState, got %v", err)
	}
}

// --------------------------------------------------
// TestRefund_CreditsBalanceBack
// --------------------------------------------------

func TestRefund_CreditsBalanceBack(t *testing.T) {
	accounts, charges := newGatewayForTest()
	account, _ := accounts.CreateAccount("Danilo", 10_000)
	charge, _ := charges.CreateCharge(account.ID, 3_000)
	_, _ = charges.Approve(charge.ID)

	refunded, err := charges.Refund(charge.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if refunded.Status != domain.ChargeStatusRefunded {
		t.Errorf("expected status refunded, got %s", refunded.Status)
	}

	updatedAccount, _ := accounts.GetAccount(account.ID)
	if updatedAccount.BalanceCents != 10_000 {
		t.Errorf("expected balance restored to 10000, got %d", updatedAccount.BalanceCents)
	}
}

// --------------------------------------------------
// TestRefund_PendingChargeCannotBeRefunded
// --------------------------------------------------

func TestRefund_PendingChargeCannotBeRefunded(t *testing.T) {
	accounts, charges := newGatewayForTest()
	account, _ := accounts.CreateAccount("Danilo", 10_000)
	charge, _ := charges.CreateCharge(account.ID, 3_000)

	_, err := charges.Refund(charge.ID)
	if err != domain.ErrInvalidChargeState {
		t.Errorf("expected ErrInvalidChargeState, got %v", err)
	}
}

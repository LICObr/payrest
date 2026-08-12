package service_test

import (
	"testing"

	"github.com/LICObr/payrest/internal/domain"
	"github.com/LICObr/payrest/internal/repository/memory"
	"github.com/LICObr/payrest/internal/service"
)

// --------------------------------------------------
// TestCreateAccount_Success
// --------------------------------------------------

func TestCreateAccount_Success(t *testing.T) {
	accounts := service.NewAccountService(memory.NewAccountRepository())

	account, err := accounts.CreateAccount("Danilo", 10_000)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if account.BalanceCents != 10_000 {
		t.Errorf("expected balance 10000, got %d", account.BalanceCents)
	}
}

// --------------------------------------------------
// TestCreateAccount_NegativeInitialBalance
// --------------------------------------------------

func TestCreateAccount_NegativeInitialBalance(t *testing.T) {
	accounts := service.NewAccountService(memory.NewAccountRepository())

	_, err := accounts.CreateAccount("Danilo", -1)
	if err != domain.ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

// --------------------------------------------------
// TestGetAccount_NotFound
// --------------------------------------------------

func TestGetAccount_NotFound(t *testing.T) {
	accounts := service.NewAccountService(memory.NewAccountRepository())

	_, err := accounts.GetAccount("does-not-exist")
	if err != domain.ErrAccountNotFound {
		t.Errorf("expected ErrAccountNotFound, got %v", err)
	}
}

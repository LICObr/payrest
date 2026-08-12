package domain

// AccountRepository persists and retrieves accounts.
// Services depend on this interface, never on a concrete storage engine.
type AccountRepository interface {
	Save(account *Account) error
	FindByID(id string) (*Account, error)
}

// ChargeRepository persists and retrieves charges.
type ChargeRepository interface {
	Save(charge *Charge) error
	FindByID(id string) (*Charge, error)
	FindByAccountID(accountID string) ([]*Charge, error)
}

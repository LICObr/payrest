# PayRest

A simulated payment gateway REST API written in Go. It models the core
lifecycle of a payment: create a charge, approve or decline it, and refund
it later — with real balance movement between accounts.

## Why this project

Payment flows are a good way to show backend fundamentals that go beyond
CRUD: a state machine, balance consistency, and domain rules that must be
enforced regardless of which client calls the API.

## Architecture

The code follows a small clean-architecture split so that business rules
never depend on HTTP or on a specific storage engine:

```
cmd/api            entry point, wires everything together
internal/domain     entities, business rules, repository interfaces (ports)
internal/service    use cases (application layer)
internal/repository/memory   in-memory implementation of the repository ports
internal/http       HTTP handlers, router, request/response DTOs
```

Swapping the in-memory storage for Postgres, for example, only requires a
new implementation of `domain.AccountRepository` / `domain.ChargeRepository`
— nothing above the repository layer changes.

## Charge state machine

```
pending ──approve──> approved ──refund──> refunded
   │
   └──decline──> declined
```

- **Approve** debits the account balance. Fails with `409` if the account
  doesn't have enough funds.
- **Decline** just closes the charge, no balance movement.
- **Refund** only works on an approved charge and credits the balance back.

## Running locally

```bash
go run ./cmd/api
```

The server starts on `:8080` (override with the `PORT` env var).

## API

| Method | Path                      | Description              |
|--------|---------------------------|---------------------------|
| POST   | `/accounts`                | Create an account         |
| GET    | `/accounts/{id}`           | Get an account             |
| POST   | `/charges`                  | Create a charge (pending)  |
| GET    | `/charges/{id}`             | Get a charge                |
| POST   | `/charges/{id}/approve`     | Approve a pending charge    |
| POST   | `/charges/{id}/decline`     | Decline a pending charge    |
| POST   | `/charges/{id}/refund`      | Refund an approved charge   |

### Example

```bash
curl -X POST localhost:8080/accounts \
  -d '{"name":"Danilo","initial_balance_cents":10000}'

curl -X POST localhost:8080/charges \
  -d '{"account_id":"<account-id>","amount_cents":3000}'

curl -X POST localhost:8080/charges/<charge-id>/approve
```

## Tests

```bash
go test ./...
```

Unit tests cover the service layer (state transitions, balance rules) and
integration tests exercise the HTTP handlers end to end with `httptest`.

## Tech stack

Go, standard library only (`net/http` with method+path routing, no external
framework) — kept intentionally small and dependency-free.

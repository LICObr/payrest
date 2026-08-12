package main

import (
	"log"
	"net/http"
	"os"

	handler "github.com/LICObr/payrest/internal/http"
	"github.com/LICObr/payrest/internal/repository/memory"
	"github.com/LICObr/payrest/internal/service"
)

func main() {
	accountRepo := memory.NewAccountRepository()
	chargeRepo := memory.NewChargeRepository()

	accountService := service.NewAccountService(accountRepo)
	chargeService := service.NewChargeService(chargeRepo, accountRepo)

	accountHandler := handler.NewAccountHandler(accountService)
	chargeHandler := handler.NewChargeHandler(chargeService)

	router := handler.NewRouter(accountHandler, chargeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("payrest listening on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

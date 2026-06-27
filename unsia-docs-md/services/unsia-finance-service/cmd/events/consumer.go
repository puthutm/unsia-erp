package main

import (
	"log"

	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/database"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	if err := infrastructure.RunEventConsumer(db); err != nil {
		log.Fatalf("Event consumer failed: %v", err)
	}
}

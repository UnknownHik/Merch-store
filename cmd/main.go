package main

import (
	"log"

	"API-Avito-shop/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("Error creating app: %v", err)
	}

	if err := a.Run(); err != nil {
		log.Fatalf("Error running app: %v", err)
	}
}

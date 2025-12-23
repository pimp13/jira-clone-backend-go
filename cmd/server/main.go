package main

import (
	"log"

	"github.com/pimp13/jira-clone-backend-go/internal/infrastructure/app"
)

func main() {
	application := app.NewApp()
	if err := application.Bootstrap(); err != nil {
		log.Fatalf("error in init application: %v", err)
	}
}

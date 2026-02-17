package main

import (
	"log"
	"os"

	_ "github.com/pimp13/jira-clone-backend-go/docs"
	"github.com/pimp13/jira-clone-backend-go/internal/infrastructure/app"
)

// @title			jira clone api
// @description	This is a jira clone api docs.
// @version		1.0
// @BasePath		/api
func main() {
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err := application.Bootstrap(); err != nil {
		log.Fatalf("bootstrap failed: %v", err)
	}

	if err := application.Start(); err != nil {
		log.Printf("server shutdown with error: %v", err)
		os.Exit(1)
	}
}

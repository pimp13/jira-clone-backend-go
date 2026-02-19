package main

import (
	"flag"
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
	var port uint
	flag.UintVar(&port, "port", 0, "Application port override")
	flag.UintVar(&port, "p", 0, "Application port override (shorthand)")
	flag.Parse()

	application, err := app.NewApp(port)
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

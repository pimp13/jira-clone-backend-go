package main

import (
	"log"

	_ "github.com/pimp13/jira-clone-backend-go/docs"
	"github.com/pimp13/jira-clone-backend-go/internal/infrastructure/app"
)

// @title			jira clone api
// @description	This is a jira clone api docs.
// @version		1.0
// @BasePath		/api
func main() {
	application := app.NewApp()
	if err := application.Bootstrap(); err != nil {
		log.Fatalf("error in init application: %v", err)
	}
}

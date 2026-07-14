package main

import (
	"log"
	"os"

	"github.com/HenryNg101/agent-service/internal/app"
	"github.com/HenryNg101/agent-service/internal/config"
)

func main() {
	log.Println("instance:", os.Getenv("HOSTNAME"))

	newApplication, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	appConfig := config.LoadAppConfig()

	router := app.SetupRouter(appConfig, newApplication.PostgresSession, newApplication)

	if err := router.Run(":" + appConfig.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

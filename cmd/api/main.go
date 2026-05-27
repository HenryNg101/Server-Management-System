package main

import (
	"log"

	"github.com/HenryNg101/server-management-system/internal/app"
	"github.com/HenryNg101/server-management-system/internal/config"
)

func main() {
	newApplication, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	appConfig := config.LoadAppConfig()

	router := app.SetupRouter(appConfig, newApplication.DB, newApplication)

	if err := router.Run(":" + appConfig.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

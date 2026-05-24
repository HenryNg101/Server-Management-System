package main

import (
	"log"

	"github.com/HenryNg101/server-management-system/internal/config"
	"github.com/HenryNg101/server-management-system/internal/database"
	"github.com/HenryNg101/server-management-system/internal/server"
)

func main() {
	postgresCfg := config.LoadPostgres()
	cfg := config.LoadAppConfig()

	postgresSession := database.NewPostgresSession(postgresCfg)

	router := server.SetupRouter(cfg, postgresSession)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

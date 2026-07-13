package postgres

import (
	"fmt"
	"log"

	"github.com/HenryNg101/cron-scheduler/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresSession(config *config.PostgresConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Host,
		config.User,
		config.Password,
		config.DatabaseName,
		config.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return db
}

// TODO: Deals with Postgres SSL mode config here
func MigratePostgres(postgresConfig *config.PostgresConfig) {
	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			postgresConfig.User,
			postgresConfig.Password,
			postgresConfig.Host,
			postgresConfig.Port,
			postgresConfig.DatabaseName,
		),
	)
	if err != nil {
		log.Fatalf("Postgres migration init error: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Postgres migration failed: %v", err)
	}

	log.Println("Postgres migrations applied (or already up-to-date)")
}

package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	"spy_cats_agency/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var direction string
	flag.StringVar(&direction, "direction", "up", "Migration direction: up or down")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load env")
		return
	}

	cfg, err := config.LoadConfig(".")
	if err != nil {
		fmt.Println("Failed to load config", slog.Any("error", err))
		os.Exit(1)
	}
	// Construct database URL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Create migration instance
	m, err := migrate.New(
		"file://db/migration",
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}
	defer m.Close()

	// Run migration based on direction
	switch direction {
	case "up":
		fmt.Println("Running migrations up...")
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				fmt.Println("No migrations to run")
				return
			}
			log.Fatalf("Failed to run migrations up: %v", err)
		}
		fmt.Println("Migrations completed successfully")
	case "down":
		fmt.Println("Running migrations down...")
		if err := m.Down(); err != nil {
			if err == migrate.ErrNoChange {
				fmt.Println("No migrations to rollback")
				return
			}
			log.Fatalf("Failed to run migrations down: %v", err)
		}
		fmt.Println("Migrations rolled back successfully")
	default:
		log.Fatalf("Invalid direction: %s. Use 'up' or 'down'", direction)
	}
}

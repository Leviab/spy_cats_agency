package postgres

import (
	"context"
	"fmt"
)

// RunMigrations executes the database migrations.
func (db *DB) RunMigrations(ctx context.Context) error {
	// Create migrations table if it doesn't exist
	createMigrationsTable := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.ExecContext(ctx, createMigrationsTable); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Check if migration already applied
	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = $1", "000001_initial_schema").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	if count > 0 {
		return nil // Migration already applied
	}

	// Run the initial schema migration
	initialSchema := `
	CREATE TABLE "cats" (
	  "id" bigserial PRIMARY KEY,
	  "name" varchar NOT NULL,
	  "years_of_experience" int NOT NULL,
	  "breed" varchar NOT NULL,
	  "salary" decimal NOT NULL,
	  "status" varchar NOT NULL DEFAULT 'available',
	  "created_at" timestamptz NOT NULL DEFAULT (now()),
	  "updated_at" timestamptz NOT NULL DEFAULT (now())
	);

	CREATE TABLE "missions" (
	  "id" bigserial PRIMARY KEY,
	  "cat_id" bigint UNIQUE,
	  "completed" boolean NOT NULL DEFAULT false,
	  "created_at" timestamptz NOT NULL DEFAULT (now()),
	  "updated_at" timestamptz NOT NULL DEFAULT (now())
	);

	CREATE TABLE "targets" (
	  "id" bigserial PRIMARY KEY,
	  "mission_id" bigint NOT NULL,
	  "name" varchar NOT NULL,
	  "country" varchar NOT NULL,
	  "notes" text NOT NULL DEFAULT '',
	  "completed" boolean NOT NULL DEFAULT false,
	  "created_at" timestamptz NOT NULL DEFAULT (now()),
	  "updated_at" timestamptz NOT NULL DEFAULT (now())
	);

	ALTER TABLE "missions" ADD FOREIGN KEY ("cat_id") REFERENCES "cats" ("id") ON DELETE SET NULL;
	ALTER TABLE "targets" ADD FOREIGN KEY ("mission_id") REFERENCES "missions" ("id") ON DELETE CASCADE;

	CREATE INDEX ON "cats" ("status");
	CREATE INDEX ON "missions" ("cat_id");
	CREATE INDEX ON "targets" ("mission_id");

	INSERT INTO schema_migrations (version) VALUES ('000001_initial_schema');
	`

	if _, err := db.ExecContext(ctx, initialSchema); err != nil {
		return fmt.Errorf("failed to run initial schema migration: %w", err)
	}

	return nil
}

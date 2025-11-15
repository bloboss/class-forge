package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

// MigrateConfig holds configuration for database migrations
type MigrateConfig struct {
	MigrationsPath string
	DatabaseName   string
}

// RunMigrations runs all pending database migrations
func RunMigrations(db *sql.DB, cfg MigrateConfig, logger *zap.Logger) error {
	if db == nil {
		return fmt.Errorf("database connection cannot be nil")
	}
	if cfg.MigrationsPath == "" {
		return fmt.Errorf("migrations path cannot be empty")
	}
	if cfg.DatabaseName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	logger.Info("Starting database migrations",
		zap.String("migrations_path", cfg.MigrationsPath),
		zap.String("database", cfg.DatabaseName),
	)

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{
		DatabaseName: cfg.DatabaseName,
	})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.MigrationsPath),
		cfg.DatabaseName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Run migrations
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("No pending migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

// RollbackMigration rolls back the last migration
func RollbackMigration(db *sql.DB, cfg MigrateConfig, logger *zap.Logger) error {
	if db == nil {
		return fmt.Errorf("database connection cannot be nil")
	}
	if cfg.MigrationsPath == "" {
		return fmt.Errorf("migrations path cannot be empty")
	}
	if cfg.DatabaseName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	logger.Info("Rolling back last migration",
		zap.String("migrations_path", cfg.MigrationsPath),
		zap.String("database", cfg.DatabaseName),
	)

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{
		DatabaseName: cfg.DatabaseName,
	})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.MigrationsPath),
		cfg.DatabaseName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Rollback one step
	if err := m.Steps(-1); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	logger.Info("Migration rolled back successfully")
	return nil
}

// GetMigrationVersion returns the current migration version
func GetMigrationVersion(db *sql.DB, cfg MigrateConfig, logger *zap.Logger) (uint, bool, error) {
	if db == nil {
		return 0, false, fmt.Errorf("database connection cannot be nil")
	}
	if cfg.MigrationsPath == "" {
		return 0, false, fmt.Errorf("migrations path cannot be empty")
	}
	if cfg.DatabaseName == "" {
		return 0, false, fmt.Errorf("database name cannot be empty")
	}

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{
		DatabaseName: cfg.DatabaseName,
	})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.MigrationsPath),
		cfg.DatabaseName,
		driver,
	)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Get current version
	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			logger.Info("No migrations have been applied yet")
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	logger.Info("Current migration version",
		zap.Uint("version", version),
		zap.Bool("dirty", dirty),
	)

	return version, dirty, nil
}

// MigrateTo migrates to a specific version
func MigrateTo(db *sql.DB, cfg MigrateConfig, version uint, logger *zap.Logger) error {
	if db == nil {
		return fmt.Errorf("database connection cannot be nil")
	}
	if cfg.MigrationsPath == "" {
		return fmt.Errorf("migrations path cannot be empty")
	}
	if cfg.DatabaseName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	logger.Info("Migrating to specific version",
		zap.String("migrations_path", cfg.MigrationsPath),
		zap.String("database", cfg.DatabaseName),
		zap.Uint("target_version", version),
	)

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{
		DatabaseName: cfg.DatabaseName,
	})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.MigrationsPath),
		cfg.DatabaseName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Migrate to specific version
	if err := m.Migrate(version); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("Already at target version")
			return nil
		}
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}

	logger.Info("Migration to target version completed successfully",
		zap.Uint("version", version),
	)
	return nil
}

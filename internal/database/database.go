package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"go.uber.org/zap"

	"code.forgejo.org/forgejo/classroom/internal/config"
)

// DB wraps sql.DB with additional functionality
type DB struct {
	*sql.DB
	logger *zap.Logger
}

// New creates a new database connection with connection pooling
func New(cfg *config.Config, logger *zap.Logger) (*DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	// Build connection string from config
	dsn := cfg.GetDatabaseDSN()

	logger.Info("Connecting to database",
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
		zap.String("database", cfg.Database.Name),
		zap.String("user", cfg.Database.User),
		zap.String("ssl_mode", cfg.Database.SSLMode),
	)

	// Open database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.Database.MaxConnections)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.Database.ConnectionMaxLifetime)

	// Verify connection with ping and timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established successfully",
		zap.Int("max_open_conns", cfg.Database.MaxConnections),
		zap.Int("max_idle_conns", cfg.Database.MaxIdleConnections),
		zap.Duration("conn_max_lifetime", cfg.Database.ConnectionMaxLifetime),
	)

	return &DB{
		DB:     db,
		logger: logger,
	}, nil
}

// Close closes the database connection gracefully
func (db *DB) Close() error {
	db.logger.Info("Closing database connection")
	if err := db.DB.Close(); err != nil {
		db.logger.Error("Failed to close database connection", zap.Error(err))
		return fmt.Errorf("failed to close database: %w", err)
	}
	db.logger.Info("Database connection closed successfully")
	return nil
}

// Ping checks if the database connection is still alive
func (db *DB) Ping(ctx context.Context) error {
	if err := db.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	return nil
}

// Stats returns database statistics
func (db *DB) Stats() sql.DBStats {
	return db.DB.Stats()
}

// HealthCheck performs a health check on the database
func (db *DB) HealthCheck(ctx context.Context) error {
	// Ping the database
	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	// Check if we can execute a simple query
	var result int
	if err := db.QueryRowContext(ctx, "SELECT 1").Scan(&result); err != nil {
		return fmt.Errorf("health check query failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("health check returned unexpected result: %d", result)
	}

	return nil
}

// WithTransaction executes a function within a database transaction
// If the function returns an error, the transaction is rolled back
// Otherwise, the transaction is committed
func (db *DB) WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	// Start transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure transaction is finalized
	defer func() {
		if p := recover(); p != nil {
			// Rollback on panic
			if rbErr := tx.Rollback(); rbErr != nil {
				db.logger.Error("Failed to rollback transaction after panic",
					zap.Any("panic", p),
					zap.Error(rbErr),
				)
			}
			panic(p) // Re-throw panic after rollback
		}
	}()

	// Execute function
	if err := fn(tx); err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(); rbErr != nil {
			db.logger.Error("Failed to rollback transaction",
				zap.Error(err),
				zap.Error(rbErr),
			)
			return fmt.Errorf("transaction failed: %w (rollback error: %v)", err, rbErr)
		}
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

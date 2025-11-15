package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"code.forgejo.org/forgejo/classroom/internal/config"
)

// getTestConfig returns a configuration for testing
func getTestConfig() *config.Config {
	return &config.Config{
		Database: config.DatabaseConfig{
			Host:                  getEnv("FGC_DATABASE_HOST", "localhost"),
			Port:                  getEnvInt("FGC_DATABASE_PORT", 5432),
			User:                  getEnv("FGC_DATABASE_USER", "fgc_test"),
			Password:              getEnv("FGC_DATABASE_PASSWORD", "fgc_test_password"),
			Name:                  getEnv("FGC_DATABASE_NAME", "forgejo_classroom_test"),
			SSLMode:               getEnv("FGC_DATABASE_SSL_MODE", "disable"),
			MaxConnections:        10,
			MaxIdleConnections:    5,
			ConnectionMaxLifetime: time.Hour,
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an integer environment variable or returns a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func TestNew(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	cfg := getTestConfig()

	t.Run("successful connection", func(t *testing.T) {
		db, err := New(cfg, logger)
		require.NoError(t, err)
		require.NotNil(t, db)
		defer db.Close()

		// Verify connection is working
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = db.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("nil config", func(t *testing.T) {
		db, err := New(nil, logger)
		assert.Error(t, err)
		assert.Nil(t, db)
		assert.Contains(t, err.Error(), "config cannot be nil")
	})

	t.Run("nil logger", func(t *testing.T) {
		db, err := New(cfg, nil)
		assert.Error(t, err)
		assert.Nil(t, db)
		assert.Contains(t, err.Error(), "logger cannot be nil")
	})

	t.Run("invalid host", func(t *testing.T) {
		invalidCfg := *cfg
		invalidCfg.Database.Host = "invalid-host-that-does-not-exist"

		db, err := New(&invalidCfg, logger)
		assert.Error(t, err)
		assert.Nil(t, db)
	})
}

func TestDB_Ping(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	cfg := getTestConfig()
	db, err := New(cfg, logger)
	require.NoError(t, err)
	defer db.Close()

	t.Run("successful ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := db.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("ping with canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		err := db.Ping(ctx)
		assert.Error(t, err)
	})
}

func TestDB_HealthCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	cfg := getTestConfig()
	db, err := New(cfg, logger)
	require.NoError(t, err)
	defer db.Close()

	t.Run("successful health check", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := db.HealthCheck(ctx)
		assert.NoError(t, err)
	})

	t.Run("health check with canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		err := db.HealthCheck(ctx)
		assert.Error(t, err)
	})
}

func TestDB_Stats(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	cfg := getTestConfig()
	db, err := New(cfg, logger)
	require.NoError(t, err)
	defer db.Close()

	t.Run("get stats", func(t *testing.T) {
		stats := db.Stats()
		assert.GreaterOrEqual(t, stats.MaxOpenConnections, 0)
		assert.GreaterOrEqual(t, stats.OpenConnections, 0)
		assert.GreaterOrEqual(t, stats.InUse, 0)
		assert.GreaterOrEqual(t, stats.Idle, 0)
	})
}

func TestDB_WithTransaction(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	cfg := getTestConfig()
	db, err := New(cfg, logger)
	require.NoError(t, err)
	defer db.Close()

	// Create a test table
	ctx := context.Background()
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS test_transactions (
			id SERIAL PRIMARY KEY,
			value TEXT NOT NULL
		)
	`)
	require.NoError(t, err)
	defer db.ExecContext(ctx, "DROP TABLE IF EXISTS test_transactions")

	t.Run("successful transaction", func(t *testing.T) {
		err := db.WithTransaction(ctx, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, "INSERT INTO test_transactions (value) VALUES ($1)", "test1")
			return err
		})
		assert.NoError(t, err)

		// Verify data was committed
		var count int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_transactions WHERE value = $1", "test1").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)

		// Cleanup
		_, _ = db.ExecContext(ctx, "DELETE FROM test_transactions WHERE value = $1", "test1")
	})

	t.Run("transaction rollback on error", func(t *testing.T) {
		testValue := "test2"
		err := db.WithTransaction(ctx, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, "INSERT INTO test_transactions (value) VALUES ($1)", testValue)
			if err != nil {
				return err
			}
			// Return an error to trigger rollback
			return assert.AnError
		})
		assert.Error(t, err)

		// Verify data was not committed
		var count int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_transactions WHERE value = $1", testValue).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("multiple operations in transaction", func(t *testing.T) {
		err := db.WithTransaction(ctx, func(tx *sql.Tx) error {
			for i := 0; i < 5; i++ {
				_, err := tx.ExecContext(ctx, "INSERT INTO test_transactions (value) VALUES ($1)", fmt.Sprintf("test_multi_%d", i))
				if err != nil {
					return err
				}
			}
			return nil
		})
		assert.NoError(t, err)

		// Verify all data was committed
		var count int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_transactions WHERE value LIKE 'test_multi_%'").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 5, count)

		// Cleanup
		_, _ = db.ExecContext(ctx, "DELETE FROM test_transactions WHERE value LIKE 'test_multi_%'")
	})
}

func TestDB_Close(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	cfg := getTestConfig()
	db, err := New(cfg, logger)
	require.NoError(t, err)

	t.Run("successful close", func(t *testing.T) {
		err := db.Close()
		assert.NoError(t, err)
	})

	t.Run("ping after close should fail", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := db.Ping(ctx)
		assert.Error(t, err)
	})
}

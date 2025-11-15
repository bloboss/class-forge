package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"code.forgejo.org/forgejo/classroom/internal/api"
	"code.forgejo.org/forgejo/classroom/internal/config"
	"code.forgejo.org/forgejo/classroom/internal/database"
	"code.forgejo.org/forgejo/classroom/internal/forgejo"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Initialize configuration
	if err := initConfig(); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Initialize logger
	logger, err := initLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting Forgejo Classroom Server",
		zap.String("version", version),
		zap.String("commit", commit),
		zap.String("build_date", date),
	)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize database connection
	db, err := database.New(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()

	logger.Info("Database initialized successfully")

	// Run database migrations
	migrateConfig := database.MigrateConfig{
		MigrationsPath: "./migrations",
		DatabaseName:   cfg.Database.Name,
	}
	if err := database.RunMigrations(db.DB, migrateConfig, logger); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err))
	}

	// Get current migration version
	version, dirty, err := database.GetMigrationVersion(db.DB, migrateConfig, logger)
	if err != nil {
		logger.Warn("Failed to get migration version", zap.Error(err))
	} else {
		logger.Info("Database migration status",
			zap.Uint("version", version),
			zap.Bool("dirty", dirty),
		)
	}

	// Initialize Forgejo client
	forgejoClient, err := initForgejoClient(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Forgejo client", zap.Error(err))
	}
	defer forgejoClient.Close()

	logger.Info("Forgejo client initialized successfully")

	// Perform startup health check on Forgejo
	logger.Info("Performing Forgejo connectivity check...")
	ctx := context.Background()
	if err := forgejoClient.HealthCheck(ctx); err != nil {
		logger.Fatal("Forgejo health check failed", zap.Error(err))
	}
	logger.Info("Forgejo connectivity verified")

	// TODO: Initialize cache
	// TODO: Initialize services

	// Initialize Gin router
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := api.NewRouter(cfg, logger)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Starting HTTP server", zap.Int("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

func initConfig() error {
	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)

	// Environment variables
	viper.SetEnvPrefix("FGC")
	viper.AutomaticEnv()

	// Config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/forgejo-classroom/")
	viper.AddConfigPath("$HOME/.config/forgejo-classroom/")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return nil
}

func initLogger() (*zap.Logger, error) {
	var config zap.Config

	if viper.GetString("server.mode") == "release" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	return config.Build()
}

func initForgejoClient(cfg *config.Config, logger *zap.Logger) (*forgejo.Client, error) {
	clientConfig := forgejo.ClientConfig{
		BaseURL:   cfg.Forgejo.BaseURL,
		Token:     cfg.Forgejo.Token,
		Timeout:   cfg.Forgejo.Timeout,
		Logger:    logger,
		UserAgent: fmt.Sprintf("forgejo-classroom/%s", version),
	}

	client, err := forgejo.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Forgejo client: %w", err)
	}

	return client, nil
}

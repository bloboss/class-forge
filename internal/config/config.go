package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Forgejo  ForgejoConfig  `mapstructure:"forgejo"`
	Cache    CacheConfig    `mapstructure:"cache"`
	Queue    QueueConfig    `mapstructure:"queue"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Host         string `mapstructure:"host"`
	Mode         string `mapstructure:"mode"` // debug, release
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	TrustedProxies []string `mapstructure:"trusted_proxies"`
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
	MaxConnections int `mapstructure:"max_connections"`
	MaxIdleConnections int `mapstructure:"max_idle_connections"`
	ConnectionMaxLifetime time.Duration `mapstructure:"connection_max_lifetime"`
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Password string        `mapstructure:"password"`
	Database int           `mapstructure:"database"`
	PoolSize int           `mapstructure:"pool_size"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

// ForgejoConfig holds Forgejo integration configuration
type ForgejoConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Token   string `mapstructure:"token"`
	Timeout time.Duration `mapstructure:"timeout"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int `mapstructure:"requests_per_minute"`
	BurstSize         int `mapstructure:"burst_size"`
}

// CacheConfig holds caching configuration
type CacheConfig struct {
	DefaultTTL        time.Duration `mapstructure:"default_ttl"`
	ClassroomTTL      time.Duration `mapstructure:"classroom_ttl"`
	AssignmentTTL     time.Duration `mapstructure:"assignment_ttl"`
	RosterTTL         time.Duration `mapstructure:"roster_ttl"`
	SubmissionTTL     time.Duration `mapstructure:"submission_ttl"`
	EnableInMemoryFallback bool     `mapstructure:"enable_in_memory_fallback"`
}

// QueueConfig holds async queue configuration
type QueueConfig struct {
	WorkerCount      int           `mapstructure:"worker_count"`
	ProcessingTimeout time.Duration `mapstructure:"processing_timeout"`
	RetryAttempts    int           `mapstructure:"retry_attempts"`
	RetryDelay       time.Duration `mapstructure:"retry_delay"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	TokenExpiration time.Duration `mapstructure:"token_expiration"`
	JWTSecret       string        `mapstructure:"jwt_secret"`
	RequireHTTPS    bool          `mapstructure:"require_https"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `mapstructure:"level"` // debug, info, warn, error
	Format     string `mapstructure:"format"` // json, console
	OutputPath string `mapstructure:"output_path"`
}

// Load loads configuration from various sources
func Load() (*Config, error) {
	config := &Config{}

	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Apply defaults
	setDefaults(config)

	// Validate configuration
	if err := validate(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// setDefaults applies default values to configuration
func setDefaults(config *Config) {
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Mode == "" {
		config.Server.Mode = "debug"
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 30
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 30
	}

	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}
	if config.Database.Port == 0 {
		config.Database.Port = 5432
	}
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "prefer"
	}
	if config.Database.MaxConnections == 0 {
		config.Database.MaxConnections = 25
	}
	if config.Database.MaxIdleConnections == 0 {
		config.Database.MaxIdleConnections = 5
	}
	if config.Database.ConnectionMaxLifetime == 0 {
		config.Database.ConnectionMaxLifetime = time.Hour
	}

	if config.Redis.Host == "" {
		config.Redis.Host = "localhost"
	}
	if config.Redis.Port == 0 {
		config.Redis.Port = 6379
	}
	if config.Redis.PoolSize == 0 {
		config.Redis.PoolSize = 10
	}
	if config.Redis.Timeout == 0 {
		config.Redis.Timeout = 5 * time.Second
	}

	if config.Forgejo.Timeout == 0 {
		config.Forgejo.Timeout = 30 * time.Second
	}
	if config.Forgejo.RateLimit.RequestsPerMinute == 0 {
		config.Forgejo.RateLimit.RequestsPerMinute = 60
	}
	if config.Forgejo.RateLimit.BurstSize == 0 {
		config.Forgejo.RateLimit.BurstSize = 10
	}

	if config.Cache.DefaultTTL == 0 {
		config.Cache.DefaultTTL = 15 * time.Minute
	}
	if config.Cache.ClassroomTTL == 0 {
		config.Cache.ClassroomTTL = 30 * time.Minute
	}
	if config.Cache.AssignmentTTL == 0 {
		config.Cache.AssignmentTTL = 15 * time.Minute
	}
	if config.Cache.RosterTTL == 0 {
		config.Cache.RosterTTL = 5 * time.Minute
	}
	if config.Cache.SubmissionTTL == 0 {
		config.Cache.SubmissionTTL = 2 * time.Minute
	}

	if config.Queue.WorkerCount == 0 {
		config.Queue.WorkerCount = 3
	}
	if config.Queue.ProcessingTimeout == 0 {
		config.Queue.ProcessingTimeout = 10 * time.Minute
	}
	if config.Queue.RetryAttempts == 0 {
		config.Queue.RetryAttempts = 3
	}
	if config.Queue.RetryDelay == 0 {
		config.Queue.RetryDelay = 30 * time.Second
	}

	if config.Auth.TokenExpiration == 0 {
		config.Auth.TokenExpiration = 24 * time.Hour
	}

	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "console"
	}
}

// validate validates the configuration
func validate(config *Config) error {
	if config.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if config.Forgejo.BaseURL == "" {
		return fmt.Errorf("forgejo base URL is required")
	}
	if config.Forgejo.Token == "" {
		return fmt.Errorf("forgejo token is required")
	}
	if config.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	// Validate enum values
	validModes := map[string]bool{"debug": true, "release": true}
	if !validModes[config.Server.Mode] {
		return fmt.Errorf("invalid server mode: %s", config.Server.Mode)
	}

	validSSLModes := map[string]bool{"disable": true, "require": true, "verify-ca": true, "verify-full": true, "prefer": true}
	if !validSSLModes[config.Database.SSLMode] {
		return fmt.Errorf("invalid database SSL mode: %s", config.Database.SSLMode)
	}

	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[config.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", config.Logging.Level)
	}

	validLogFormats := map[string]bool{"json": true, "console": true}
	if !validLogFormats[config.Logging.Format] {
		return fmt.Errorf("invalid log format: %s", config.Logging.Format)
	}

	return nil
}

// GetDatabaseDSN returns the PostgreSQL connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:"server" yaml:"server"`
	Database DatabaseConfig `mapstructure:"database" yaml:"database"`
	Auth     AuthConfig     `mapstructure:"auth" yaml:"auth"`
	Redis    RedisConfig    `mapstructure:"redis" yaml:"redis"`
	CORS     CORSConfig     `mapstructure:"cors" yaml:"cors"`
	Log      LogConfig      `mapstructure:"log" yaml:"log"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Host    string        `mapstructure:"host" yaml:"host"`
	Port    int           `mapstructure:"port" yaml:"port"`
	Timeout TimeoutConfig `mapstructure:"timeout" yaml:"timeout"`
}

// TimeoutConfig holds timeout-related configuration
type TimeoutConfig struct {
	Server string `mapstructure:"server" yaml:"server"`
	Write  string `mapstructure:"write" yaml:"write"`
	Read   string `mapstructure:"read" yaml:"read"`
	Idle   string `mapstructure:"idle" yaml:"idle"`
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Driver             string `mapstructure:"driver" yaml:"driver"`
	Host               string `mapstructure:"host" yaml:"host"`
	Port               int    `mapstructure:"port" yaml:"port"`
	Name               string `mapstructure:"name" yaml:"name"`
	Schema             string `mapstructure:"schema" yaml:"schema"`
	User               string `mapstructure:"user" yaml:"user"`
	Password           string `mapstructure:"password" yaml:"password"`
	SSLMode            string `mapstructure:"sslmode" yaml:"sslmode"`
	MaxConnections     int    `mapstructure:"max_connections" yaml:"max_connections"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections" yaml:"max_idle_connections"`
	ConnMaxLifetime    string `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	StatementTimeout   string `mapstructure:"statement_timeout" yaml:"statement_timeout"`
	LockTimeout        string `mapstructure:"lock_timeout" yaml:"lock_timeout"`
	Pool               struct {
		MaxConns    int    `mapstructure:"max_conns" yaml:"max_conns"`
		MinConns    int    `mapstructure:"min_conns" yaml:"min_conns"`
		MaxLifetime string `mapstructure:"max_lifetime" yaml:"max_lifetime"`
		MaxIdleTime string `mapstructure:"max_idle_time" yaml:"max_idle_time"`
	} `mapstructure:"pool" yaml:"pool"`
}

// AuthConfig holds authentication-specific configuration
type AuthConfig struct {
	ServiceURL string `mapstructure:"service_url" yaml:"service_url"`
	Username   string `mapstructure:"username" yaml:"username"`
	Password   string `mapstructure:"password" yaml:"password"`
	Timeout    string `mapstructure:"timeout" yaml:"timeout"`
}

// RedisConfig holds Redis-specific configuration
type RedisConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Password string `mapstructure:"password" yaml:"password"`
	DB       int    `mapstructure:"db" yaml:"db"`
}

// CORSConfig holds CORS-specific configuration
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins" yaml:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods" yaml:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers" yaml:"allowed_headers"`
}

// LogConfig holds logging-specific configuration
type LogConfig struct {
	Level  string `mapstructure:"level" yaml:"level"`
	Format string `mapstructure:"format" yaml:"format"`
}

// Validate performs validation on the configuration
func (c *Config) Validate() error {
	// Validate Server configuration
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	// Validate timeouts
	timeoutFields := map[string]string{
		"server timeout": c.Server.Timeout.Server,
		"write timeout":  c.Server.Timeout.Write,
		"read timeout":   c.Server.Timeout.Read,
		"idle timeout":   c.Server.Timeout.Idle,
	}

	for name, value := range timeoutFields {
		if _, err := time.ParseDuration(value); err != nil {
			return fmt.Errorf("invalid %s: %s", name, value)
		}
	}

	// Validate Database configuration
	if err := c.validateDatabaseConfig(); err != nil {
		return fmt.Errorf("database configuration error: %w", err)
	}

	return nil
}

func (c *Config) validateDatabaseConfig() error {
	switch c.Database.Driver {
	case "postgres":
		if c.Database.Host == "" {
			return fmt.Errorf("database host is required for PostgreSQL")
		}
		if c.Database.Port <= 0 || c.Database.Port > 65535 {
			return fmt.Errorf("invalid database port: %d", c.Database.Port)
		}
		if c.Database.Name == "" {
			return fmt.Errorf("database name is required")
		}
		if c.Database.User == "" {
			return fmt.Errorf("database user is required")
		}
	case "mysql":
		if c.Database.Host == "" {
			return fmt.Errorf("database host is required for MySQL")
		}
		if c.Database.Port <= 0 || c.Database.Port > 65535 {
			return fmt.Errorf("invalid database port: %d", c.Database.Port)
		}
	case "sqlite3":
		if c.Database.Name == "" {
			return fmt.Errorf("database name (file path) is required for SQLite")
		}
	default:
		return fmt.Errorf("unsupported database driver: %s", c.Database.Driver)
	}

	// Validate connection pool settings
	if c.Database.MaxConnections < 1 {
		return fmt.Errorf("max_connections must be at least 1")
	}
	if c.Database.MaxIdleConnections > c.Database.MaxConnections {
		return fmt.Errorf("max_idle_connections cannot be greater than max_connections")
	}
	if c.Database.ConnMaxLifetime != "" {
		if _, err := time.ParseDuration(c.Database.ConnMaxLifetime); err != nil {
			return fmt.Errorf("invalid connection max lifetime: %s", c.Database.ConnMaxLifetime)
		}
	}

	return nil
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (*Config, error) {
	// First try to load .env file
	if err := godotenv.Load("config/.env"); err != nil {
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	// Add environment variable support
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables
	envMappings := map[string]string{
		"database.host":     "DB_HOST",
		"database.port":     "DB_PORT",
		"database.name":     "DB_NAME",
		"database.user":     "DB_USER",
		"database.password": "DB_PASSWORD",
		"database.sslmode":  "DB_SSLMODE",
		"server.port":       "SERVER_PORT",
		"log.level":         "LOG_LEVEL",
	}

	for configKey, envVar := range envMappings {
		viper.BindEnv(configKey, "APP_"+envVar)
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set defaults if not provided
	setDefaults(&config)

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func setDefaults(config *Config) {
	if config.Database.MaxConnections == 0 {
		config.Database.MaxConnections = 20
	}
	if config.Database.MaxIdleConnections == 0 {
		config.Database.MaxIdleConnections = 5
	}
	if config.Database.ConnMaxLifetime == "" {
		config.Database.ConnMaxLifetime = "1h"
	}
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}
	if config.Log.Level == "" {
		config.Log.Level = "info"
	}
	if config.Log.Format == "" {
		config.Log.Format = "json"
	}
}

// NewDatabase creates a new database connection
func (c *Config) NewDatabase() (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch c.Database.Driver {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Database.Host,
			c.Database.Port,
			c.Database.User,
			c.Database.Password,
			c.Database.Name,
			c.Database.SSLMode,
		)
		if c.Database.Schema != "" {
			dsn += fmt.Sprintf(" search_path=%s", c.Database.Schema)
		}
		dialector = postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		})
	case "mysql":
		dialector = mysql.Open(c.Database.buildMySQLDSN())
	case "sqlite3":
		dialector = sqlite.Open(c.Database.Name)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", c.Database.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(c.Database.MaxConnections)
	sqlDB.SetMaxIdleConns(c.Database.MaxIdleConnections)
	if c.Database.ConnMaxLifetime != "" {
		duration, err := time.ParseDuration(c.Database.ConnMaxLifetime)
		if err != nil {
			return nil, fmt.Errorf("invalid connection max lifetime: %w", err)
		}
		sqlDB.SetConnMaxLifetime(duration)
	}

	return db, nil
}

// buildMySQLDSN builds the DSN string for MySQL
func (c *DatabaseConfig) buildMySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
}

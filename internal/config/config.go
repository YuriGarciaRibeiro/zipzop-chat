package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
	SSLMode      string
	MaxConns     int
	ConnTimeout  time.Duration
}

type AuthConfig struct {
	JWTSecret    string
	JWTExpiry    time.Duration
	PasswordSalt string
	TokenIssuer  string
}

type ServerConfig struct {
	Port         string
	Debug        bool
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type AppConfig struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

func LoadConfig() (*AppConfig, error) {
	if err := Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	// Carrega configuração do banco
	dbConfig, err := loadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("database config error: %w", err)
	}

	// Carrega configuração de autenticação
	authConfig, err := loadAuthConfig()
	if err != nil {
		return nil, fmt.Errorf("auth config error: %w", err)
	}

	// Carrega configuração do servidor
	serverConfig := loadServerConfig()

	return &AppConfig{
		Server:   serverConfig,
		Database: *dbConfig,
		Auth:     *authConfig,
	}, nil
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port:         getEnv("PORT", "8080"),
		Debug:        getEnvBool("DEBUG", false),
		ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
		WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
		IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
	}
}

func loadDatabaseConfig() (*DatabaseConfig, error) {
	cfg := &DatabaseConfig{
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         getEnv("DB_PORT", "5432"),
		User:         getEnv("DB_USER", ""),
		Password:     getEnv("DB_PASSWORD", ""),
		DatabaseName: getEnv("DB_NAME", "zipzop"),
		SSLMode:      getEnv("DB_SSL_MODE", "disable"),
		MaxConns:     getEnvInt("DB_MAX_CONNS", 20),
		ConnTimeout:  getEnvDuration("DB_CONN_TIMEOUT", 10*time.Second),
	}

	if cfg.User == "" || cfg.Password == "" {
		return nil, fmt.Errorf("database credentials not provided")
	}

	return cfg, nil
}

func loadAuthConfig() (*AuthConfig, error) {
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET not provided")
	}

	return &AuthConfig{
		JWTSecret:    jwtSecret,
		JWTExpiry:    getEnvDuration("JWT_EXPIRY", 24*time.Hour),
		PasswordSalt: getEnv("PASSWORD_SALT", ""),
		TokenIssuer:  getEnv("TOKEN_ISSUER", "zipzop-chat"),
	}, nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}
	return b
}

func getEnvInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return i
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultValue
	}
	return d
}

package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	Env        string
	ServerPort string
	LogLevel   string
	Database   DatabaseConfig
	RabbitMQ   RabbitMQConfig
	MockMode   bool
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RabbitMQConfig struct {
	URL            string
	MaxRetries     int
	RetryDelay     time.Duration
	PrefetchCount  int
	ReconnectDelay time.Duration
}

func LoadConfig() (*AppConfig, error) {
	return &AppConfig{
		Env:        getEnv("ENV", "development"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "notifications"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		RabbitMQ: RabbitMQConfig{
			URL:            getEnv("RABBITMQ_URL", "amqp://admin:admin123@tokohobby-rabbitmq:5672/tokohobby"),
			MaxRetries:     getEnvInt("RABBITMQ_MAX_RETRIES", 3),
			RetryDelay:     time.Duration(getEnvInt("RABBITMQ_RETRY_DELAY", 5)) * time.Second,
			PrefetchCount:  getEnvInt("RABBITMQ_PREFETCH_COUNT", 10),
			ReconnectDelay: time.Duration(getEnvInt("RABBITMQ_RECONNECT_DELAY", 10)) * time.Second,
		},
		MockMode: getEnvBool("MOCK_MODE", true),
	}, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

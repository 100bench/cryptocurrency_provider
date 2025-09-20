package config

import (
	"fmt"
	"os"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Kafka
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string

	// External API
	CoinDeskAPIURL string
}

func Load() *Config {
	return &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "currency_db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// Kafka
		KafkaBrokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
		KafkaTopic:   getEnv("KAFKA_TOPIC", "currency-rates"),
		KafkaGroupID: getEnv("KAFKA_GROUP_ID", "currency-consumer-group"),

		// External API
		CoinDeskAPIURL: getEnv("COINDESK_API_URL", "https://api.coindesk.com/v1/bpi/currentprice.json"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// PostgresDSN возвращает DSN строку для подключения к PostgreSQL
func (c *Config) PostgresDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// Kafka
	KafkaBrokers []string `mapstructure:"kafka_brokers"`
	KafkaTopic   string   `mapstructure:"kafka_topic"`
	KafkaGroupID string   `mapstructure:"kafka_group_id"`

	// External API
	CoinDeskAPIURL string `mapstructure:"coindesk_api_url"`

	// HTTP Server
	HTTPAddr string `mapstructure:"http_addr"`

	// Fetch settings
	FetchInterval   time.Duration `mapstructure:"fetch_interval"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

var v *viper.Viper

func init() {
	v = viper.New()

	// Настройка Viper
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./deployment/config")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

}

func Load() *Config {
	// Читаем конфигурационный файл
	if err := v.ReadInConfig(); err != nil {
		// Если файл не найден, используем только переменные окружения и defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("unable to decode config into struct: %w", err))
	}

	return &config
}

// PostgresDSN возвращает DSN строку для подключения к PostgreSQL
func (c *Config) PostgresDSN() string {
	connString := v.GetString("connection_string")
	return connString
}

// GetViper возвращает экземпляр Viper для дополнительных операций
func GetViper() *viper.Viper {
	return v
}

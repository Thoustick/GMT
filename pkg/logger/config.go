package logger

import (
	"github.com/spf13/viper"
	"log"
)

// Config структура конфигурации
type Config struct {
	LogLevel string
}

// LoadConfig загружает конфигурацию из .env или переменных окружения
func LoadConfig() (*Config, error) {
	viper.SetConfigFile("config/.env") // Загружаем .env
	viper.ReadInConfig()
	// Читаем конфигурацию, но не падаем с ошибкой, если файла нет
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read .env file, using environment variables only")
	}

	// Устанавливаем значения по умолчанию
	viper.SetDefault("log.level", "info")

	config := &Config{
		LogLevel: viper.GetString("log.level"),
	}

	return config, nil
}

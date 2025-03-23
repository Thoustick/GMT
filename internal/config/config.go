package config

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

// Config структура для хранения конфигурационных параметров
type Config struct {
	TelegramToken   string `mapstructure:"telegram_token"`
	RedisURL        string `mapstructure:"redis_url"`
	HugFaceApiKey   string `mapstructure:"huggingface_api_key"`
	HugFaceModel    string `mapstructure:"huggingface_model"`
	LogLevel        string `mapstructure:"log_level"`
}

var once sync.Once
var AppConfig *Config

// LoadConfig загружает конфигурацию из файла и переменных окружения
func LoadConfig() *Config {
	tmpLog := log.New(os.Stdout, "CONFIG: ", log.LstdFlags)

	v := viper.New()

	// Указываем пути поиска
	v.AddConfigPath("config")   // Директория с конфигом
	v.SetConfigName("config")    // Имя файла (без .yaml)
	v.SetConfigType("yaml")      // Тип файла

	// Читаем конфиг-файл, если он есть
	if err := v.ReadInConfig(); err != nil {
		tmpLog.Printf("Config file not found: %v. Relying on environment variables.", err)
	}

	// Автоматически подхватываем переменные окружения
	v.AutomaticEnv()

	// Заполняем структуру Config
	AppConfig = &Config{}
	if err := v.Unmarshal(AppConfig); err != nil {
		tmpLog.Fatalf("Error unmarshaling config: %v", err)
	}

	// Проверяем, загружены ли критически важные переменные
	if AppConfig.TelegramToken == "" {
		tmpLog.Fatal("TELEGRAM_TOKEN is required but missing")
	}

	return AppConfig
}

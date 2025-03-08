package logger

import (
	"os"
	"time"


	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	// Загружаем конфигурацию
	cfg, err := LoadConfig()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load config, using default values")
	}

	// Парсим уровень логирования
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Warn().Err(err).Msgf("Invalid log level: %s, defaulting to INFO", cfg.LogLevel)
		logLevel = zerolog.InfoLevel
	}

	// Определяем writer для логов (консоль)
	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	log.Logger = zerolog.New(writer).
		Level(logLevel).
		With().Timestamp().
		Logger()

	// Устанавливаем глобальный уровень логирования
	zerolog.SetGlobalLevel(logLevel)

	log.Info().Msg("Logger initialized successfully")
}

func Debug(msg string) {
	log.Debug().Msg(msg)
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Error(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

func Fatal(err error, msg string) {
	log.Fatal().Err(err).Msg(msg)
}

package main

import (
	"os"
	"os/signal"
	"syscall"
	"context"
	"time"

	"github.com/Thoustick/GMT/internal/bot"
	"github.com/Thoustick/GMT/internal/config"
	"github.com/Thoustick/GMT/internal/huggingface"
	"github.com/Thoustick/GMT/internal/tasks"
	"github.com/Thoustick/GMT/pkg/logger"
)

func main() {
	// 1️⃣ Загружаем конфигурацию
	cfg := config.LoadConfig()
	if cfg == nil {
		panic("Ошибка загрузки конфигурации")
	}

	// 2️⃣ Создаем логгер
	log := logger.InitLogger(cfg)

	// 3️⃣ Создаем Hugging Face клиент
	hfClient, err := huggingface.NewClient(cfg, log)
	if err != nil {
		log.Fatal("Ошибка создания Hugging Face клиента", err, nil)
		os.Exit(1)
	}

	// 4️⃣ Создаем генератор задач
	taskGen := tasks.NewTaskGeneratorImpl(hfClient, log)

	// 5️⃣ Запускаем Telegram-бота
	botInstance, err := bot.NewBot(cfg, log, &taskGen)
	if err != nil {
		log.Fatal("Ошибка запуска Telegram-бота", err, nil)
		os.Exit(1)
	}

	// 6️⃣ Запускаем бота в горутине
	go func() {
		if err := botInstance.Run(); err != nil {
			log.Fatal("Ошибка во время работы бота", err, nil)
			os.Exit(1)
		}
	}()

	// 7️⃣ Обрабатываем graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Info("Recived shutdown request.", nil)

	// 8️⃣ Создаем контекст с таймаутом 10 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := botInstance.Shutdown(ctx); err != nil {
		log.Error("Ошибка при остановке бота", err, nil)
	}
}

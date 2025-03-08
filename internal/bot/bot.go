package bot

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/Thoustick/GMT/pkg/logger"
)

func Init() error {
	logger.InitLogger()

	token, err := LoadConfig() 
	if err != nil {
		logger.Fatal(err, "failed to load config")
		return err
	}

	// Создание нового бота
	bot, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{
			Client: http.Client{
				Timeout: 10* time.Second,
			},
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: gotgbot.DefaultTimeout,
				APIURL:  gotgbot.DefaultAPIURL,
			},
		},
	})
	if err != nil {
		logger.Fatal(err, "failed to create bot")
		return err
	}
	logger.Info("Bot started")

	// Создание диспетчера с обработкой ошибок
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			logger.Error(err, "error in dispatcher")
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(dispatcher, nil)

	dispatcher.AddHandler(handlers.NewMessage(message.Text, HandleMessage))

	// Запуск получения обновлений
	err = updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		logger.Fatal(err, "failed to start polling")
		return err
	}

	updater.Idle()

	logger.Info("Bot stopped")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Info("Shutting down gracefully")

	updater.Stop()
	
	return nil
}
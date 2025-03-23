package bot

import (
	"context"
	"net/http"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/Thoustick/GMT/internal/config"
	"github.com/Thoustick/GMT/internal/tasks"
	"github.com/Thoustick/GMT/pkg/logger"
)
type Bot struct {
	bot          *gotgbot.Bot
	cfg          *config.Config
	log          logger.Logger
	taskGen      tasks.TaskGenerator
	dispatcher   *ext.Dispatcher
	updater      *ext.Updater
}

func NewBot(cfg *config.Config, log logger.Logger, taskGen tasks.TaskGenerator) (*Bot, error) {
	bot, err := gotgbot.NewBot(cfg.TelegramToken, &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{
			Client: http.Client{
				Timeout: 10 * time.Second,
			},
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: gotgbot.DefaultTimeout,
				APIURL:  gotgbot.DefaultAPIURL,
			},
		},
	})
	if err != nil {
		log.Error("Failed to create bot", err, nil)
		return nil, err
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Error("Error in dispatcher", err, nil)
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(dispatcher, nil)

	bh := NewHandler(taskGen, log)

	dispatcher.AddHandler(handlers.NewCommand("start", bh.StartCommand))
	dispatcher.AddHandler(handlers.NewCommand("task", bh.TaskCommand))

	return &Bot{
		bot:        bot,
		cfg:        cfg,
		log:        log,
		taskGen:    taskGen,
		dispatcher: dispatcher,
		updater:    updater,
	}, nil
}

func (b *Bot) Run() error {
	b.log.Info("Bot started", nil)

	err := b.updater.StartPolling(b.bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		b.log.Error("Failed to start polling", err, nil)
		return err
	}

	go b.updater.Idle()

	return nil
}

func (b *Bot) Shutdown(ctx context.Context) error {
	b.log.Info("Stopping bot gracefully ...", nil)

	b.updater.Stop()

	select {
	case <-time.After(2 * time.Second):
		b.log.Info("Bot shutdown complete.", nil)
	case <-ctx.Done():
		b.log.Warn("Timeout reached, forcefully stopping bot.", nil)
	}

	return nil
}

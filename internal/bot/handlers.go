package bot

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/Thoustick/GMT/internal/tasks"
	"github.com/Thoustick/GMT/pkg/logger"
)

type BotHandler struct{
	TaskGenerator tasks.TaskGenerator
	Log logger.Logger
}

func NewHandler(tg tasks.TaskGenerator, log logger.Logger) *BotHandler{
	return &BotHandler{
		TaskGenerator: tg,
		Log: log,
	}
}
// StartCommand обрабатывает команду /start.
func (bh *BotHandler) StartCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	message := fmt.Sprintf("Привет, я @%s. Используй /task для получения задачи!", b.User.Username)
	_, err := ctx.EffectiveMessage.Reply(b, message, nil)
	if err != nil {
		bh.Log.Error("Ошибка отправки сообщения /start", err, nil)
	}
	return err
}

// TaskCommand responds to the /task command by generating and sending a Markdown-formatted task.
func (bh *BotHandler) TaskCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	markdownTask, err := bh.TaskGenerator.GenerateTask()
	if err != nil {
		bh.Log.Error("Ошибка генерации задачи", err, nil)
		_, _ = ctx.EffectiveMessage.Reply(b, "Не удалось сгенерировать задачу, попробуйте позже.", nil)
		return err
	}

	_, err = ctx.EffectiveMessage.Reply(b, markdownTask, &gotgbot.SendMessageOpts{ParseMode: "MarkdownV2"})
	return err
}

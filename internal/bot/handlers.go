package bot

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/Thoustick/GMT/pkg/logger"
)

// HandleMessage обрабатывает входящие сообщения
func HandleMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	msgText := ctx.EffectiveMessage.Text
	logger.Info(msgText)

	_, err := ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	return nil
}
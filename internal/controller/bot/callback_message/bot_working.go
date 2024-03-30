package start

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/pkg/client/telegram"
)

type CallbackBotMessage struct {
	logger *slog.Logger
	bot    telegram.Bot
	tgUser tg.TgUserDTO
}

func NewCallbackBot(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO) CallbackBotMessage {
	return CallbackBotMessage{
		logger: logger,
		bot:    bot,
		tgUser: tgUser,
	}
}

// Execute место связи telegram и бизнес логи
func (c *CallbackBotMessage) Execute(message tg.MessageDTO, callbackData string) {
	msg := tgbotapi.NewMessage(c.tgUser.TgID, fmt.Sprintf("Chosen language %s", callbackData))

	_, err := c.bot.Bot.Send(msg)
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s", err))
	}
}

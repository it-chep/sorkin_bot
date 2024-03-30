package start

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/pkg/client/telegram"
)

type StartBotCommand struct {
	logger *slog.Logger
	bot    telegram.Bot
	tgUser tg.TgUserDTO
}

func NewStartBotCommand(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO) StartBotCommand {
	return StartBotCommand{
		logger: logger,
		bot:    bot,
		tgUser: tgUser,
	}
}

// Execute место связи telegram и бизнес логи
func (c *StartBotCommand) Execute(message tg.MessageDTO) {

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 EN", "EN"),
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 RU", "RU"),
		),
	)

	msg := tgbotapi.NewMessage(c.tgUser.TgID, "Before you start using the bot, please select a language")
	msg.ReplyMarkup = keyboard

	_, err := c.bot.Bot.Send(msg)
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s", err))
	}
}

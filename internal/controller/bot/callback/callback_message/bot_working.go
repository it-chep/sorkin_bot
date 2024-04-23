package start

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/services/message"
	"sorkin_bot/internal/domain/services/user"
	"sorkin_bot/pkg/client/telegram"
)

var languagesMap = map[string]bool{
	"RU": true,
	"EN": true,
}

type CallbackBotMessage struct {
	logger         *slog.Logger
	bot            telegram.Bot
	tgUser         tg.TgUserDTO
	userService    user.UserService
	messageService message.MessageService
}

func NewCallbackBot(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO, userService user.UserService, messageService message.MessageService) CallbackBotMessage {
	return CallbackBotMessage{
		logger:         logger,
		bot:            bot,
		tgUser:         tgUser,
		userService:    userService,
		messageService: messageService,
	}
}

// Execute место связи telegram и бизнес логи
func (c *CallbackBotMessage) Execute(ctx context.Context, message tg.MessageDTO, callbackData string) {
	var msg tgbotapi.MessageConfig
	if _, ok := languagesMap[callbackData]; ok {
		_, err := c.userService.ChangeLanguage(ctx, c.tgUser, callbackData)
		if err != nil {
			return
		}
		msg = tgbotapi.NewMessage(c.tgUser.TgID, "Язык поставлен")

	} else {
		msg = tgbotapi.NewMessage(c.tgUser.TgID, fmt.Sprintf("Chosen language %s", callbackData))
	}
	_, err := c.bot.Bot.Send(msg)
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s", err))
	}
}

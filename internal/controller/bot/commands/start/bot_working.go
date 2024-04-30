package start

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/services/bot"
	"sorkin_bot/internal/domain/services/message"
	"sorkin_bot/internal/domain/services/user"
	"sorkin_bot/pkg/client/telegram"
)

type StartBotCommand struct {
	logger         *slog.Logger
	bot            telegram.Bot
	tgUser         tg.TgUserDTO
	userService    user.UserService
	messageService message.MessageService
	botService     bot.BotService
}

func NewStartBotCommand(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO, userService user.UserService, messageService message.MessageService, botService bot.BotService) StartBotCommand {
	return StartBotCommand{
		logger:         logger,
		bot:            bot,
		tgUser:         tgUser,
		userService:    userService,
		messageService: messageService,
		botService:     botService,
	}
}

// Execute место связи telegram и бизнес логи
func (c *StartBotCommand) Execute(ctx context.Context, message tg.MessageDTO) {
	user, err := c.userService.RegisterNewUser(ctx, c.tgUser)
	var msg tgbotapi.MessageConfig
	var msgText string
	var keyboard tgbotapi.InlineKeyboardMarkup
	if err != nil {
		return
	}

	if user.GetState() == "" && user.GetLanguageCode() == "" {
		msgText, keyboard = c.botService.ConfigureChangeLanguageMessage(ctx, user)
	} else {
		msgText, err = c.messageService.GetMessage(ctx, user, "Start")
	}

	msg = tgbotapi.NewMessage(c.tgUser.TgID, msgText)

	if len(keyboard.InlineKeyboard) != 0 {
		msg.ReplyMarkup = keyboard
	}

	_, err = c.bot.Bot.Send(msg)
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s", err))
	}
}

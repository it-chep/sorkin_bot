package start

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/pkg/client/telegram"
)

type StartBotCommand struct {
	logger         *slog.Logger
	bot            telegram.Bot
	tgUser         tg.TgUserDTO
	userService    userService
	messageService messageService
	botService     botService
}

func NewStartBotCommand(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO, userService userService, messageService messageService, botService botService) StartBotCommand {
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

	if user.GetState() == nil && user.GetLanguageCode() == nil {
		msgText, keyboard = c.botService.ConfigureChangeLanguageMessage(ctx, user)
	} else {
		msgText, err = c.messageService.GetMessage(ctx, user, "Start")
	}

	msg = tgbotapi.NewMessage(c.tgUser.TgID, msgText)

	if keyboard.InlineKeyboard != nil {
		msg.ReplyMarkup = keyboard
	}

	c.bot.SendMessage(msg, message)
}

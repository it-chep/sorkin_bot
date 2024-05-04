package change_language

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/pkg/client/telegram"
)

type ChangeLanguageCommand struct {
	logger         *slog.Logger
	bot            telegram.Bot
	tgUser         tg.TgUserDTO
	userService    UserService
	messageService MessageService
	botService     BotService
}

func NewChangeLanguageCommand(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO, userService UserService, messageService MessageService, botService BotService) ChangeLanguageCommand {
	return ChangeLanguageCommand{
		logger:         logger,
		bot:            bot,
		tgUser:         tgUser,
		userService:    userService,
		messageService: messageService,
		botService:     botService,
	}
}

func (c ChangeLanguageCommand) Execute(ctx context.Context) {
	userEntity, err := c.userService.GetUser(ctx, c.tgUser)
	if err != nil {
		return
	}

	msgText, keyboard := c.botService.ConfigureChangeLanguageMessage(ctx, userEntity)
	msg := tgbotapi.NewMessage(c.tgUser.TgID, msgText)

	msg.ReplyMarkup = keyboard

	_, err = c.bot.Bot.Send(msg)
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s", err))
	}

}

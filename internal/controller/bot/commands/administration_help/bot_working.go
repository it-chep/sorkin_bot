package administration_help

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/pkg/client/telegram"
)

type AdministrationHelpCommand struct {
	logger         *slog.Logger
	bot            telegram.Bot
	tgUser         tg.TgUserDTO
	messageService messageService
	userService    userService
}

func NewAdministrationHelpCommand(
	logger *slog.Logger,
	bot telegram.Bot,
	tgUser tg.TgUserDTO,
	messageService messageService,
	userService userService,
) AdministrationHelpCommand {
	return AdministrationHelpCommand{
		logger:         logger,
		bot:            bot,
		tgUser:         tgUser,
		messageService: messageService,
		userService:    userService,
	}
}

func (c AdministrationHelpCommand) Execute(ctx context.Context, tgMessage tg.MessageDTO) {
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser.TgID)
	msgText, err := c.messageService.GetMessage(ctx, userEntity, "i will call tech_support")
	if err != nil {
		return
	}
	c.bot.SendMessage(tgbotapi.NewMessage(c.tgUser.TgID, msgText), tgMessage)
}

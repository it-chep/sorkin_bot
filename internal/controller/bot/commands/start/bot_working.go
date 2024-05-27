package start

import (
	"context"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
)

type StartBotCommand struct {
	logger         *slog.Logger
	botGateway     botGateway
	tgUser         tg.TgUserDTO
	machine        *state_machine.UserStateMachine
	userService    userService
	messageService messageService
}

func NewStartBotCommand(logger *slog.Logger, botGateway botGateway, tgUser tg.TgUserDTO, machine *state_machine.UserStateMachine, userService userService, messageService messageService) StartBotCommand {
	return StartBotCommand{
		logger:         logger,
		botGateway:     botGateway,
		tgUser:         tgUser,
		machine:        machine,
		userService:    userService,
		messageService: messageService,
	}
}

// Execute место связи telegram и бизнес логи
func (c *StartBotCommand) Execute(ctx context.Context, message tg.MessageDTO) {
	user, err := c.userService.RegisterNewUser(ctx, c.tgUser)
	if err != nil {
		return
	}

	if user.GetState() == nil && user.GetLanguageCode() == nil {
		c.botGateway.SendChangeLanguageMessage(ctx, user, message)
		c.machine.SetState(user, state_machine.ChooseLanguage)
	} else {
		c.botGateway.SendStartMessage(ctx, user, message)
		c.machine.SetState(user, state_machine.Start)
	}
}

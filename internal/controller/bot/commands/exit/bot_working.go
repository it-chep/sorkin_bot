package exit

import (
	"context"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/pkg/client/telegram"
)

type ExitBotCommand struct {
	logger *slog.Logger
	bot    telegram.Bot

	botGateway       botGateway
	tgUser           tg.TgUserDTO
	machine          *state_machine.UserStateMachine
	userService      userService
	draftAppointment draftAppointment
}

func NewExitBotCommand(
	logger *slog.Logger,
	bot telegram.Bot,
	botGateway botGateway,
	tgUser tg.TgUserDTO,
	machine *state_machine.UserStateMachine,
	userService userService,
	draftAppointment draftAppointment,
) ExitBotCommand {
	return ExitBotCommand{
		logger:           logger,
		botGateway:       botGateway,
		bot:              bot,
		tgUser:           tgUser,
		machine:          machine,
		userService:      userService,
		draftAppointment: draftAppointment,
	}
}

// Execute место связи telegram и бизнес логи
func (c *ExitBotCommand) Execute(ctx context.Context, message tg.MessageDTO) {
	user, err := c.userService.GetUser(ctx, c.tgUser.TgID)
	if err != nil {
		return
	}
	c.bot.RemoveMessage(user.GetTgId(), int(message.MessageID)-1)
	c.draftAppointment.CleanDraftAppointment(ctx, user.GetTgId())
	c.botGateway.SendStartMessage(ctx, user, message)
	_, _ = c.userService.ChangeState(ctx, user.GetTgId(), state_machine.Start)
}

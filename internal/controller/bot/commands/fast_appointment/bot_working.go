package fast_appointment

import (
	"context"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/pkg/client/telegram"
)

type FastAppointmentBotCommand struct {
	logger             *slog.Logger
	bot                telegram.Bot
	tgUser             tg.TgUserDTO
	userService        userService
	machine            *state_machine.UserStateMachine
	appointmentService appointmentService
	botGateway         botGateway
}

func NewFastAppointmentBotCommand(
	logger *slog.Logger,
	bot telegram.Bot,
	tgUser tg.TgUserDTO,
	userService userService,
	machine *state_machine.UserStateMachine,
	appointmentService appointmentService,
	botGateway botGateway,
) FastAppointmentBotCommand {
	return FastAppointmentBotCommand{
		logger:             logger,
		bot:                bot,
		tgUser:             tgUser,
		userService:        userService,
		machine:            machine,
		botGateway:         botGateway,
		appointmentService: appointmentService,
	}
}

func (c *FastAppointmentBotCommand) Execute(ctx context.Context, message tg.MessageDTO) {
	userEntity, err := c.userService.GetUser(ctx, c.tgUser.TgID)
	if err != nil {
		return
	}

	c.botGateway.SendFastAppointmentMessage(ctx, userEntity, message)
	c.machine.SetState(userEntity, state_machine.FastAppointment)
	c.appointmentService.CreateDraftAppointment(ctx, userEntity.GetTgId())
}

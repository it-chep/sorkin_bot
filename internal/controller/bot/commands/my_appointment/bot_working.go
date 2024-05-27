package my_appointment

import (
	"context"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
)

type MyAppointmentsCommand struct {
	logger             *slog.Logger
	botGateway         botGateway
	tgUser             tg.TgUserDTO
	machine            *state_machine.UserStateMachine
	userService        userService
	appointmentService appointmentService
}

func NewMyAppointmentsCommand(logger *slog.Logger, botGateway botGateway, tgUser tg.TgUserDTO, machine *state_machine.UserStateMachine, userService userService, appointmentService appointmentService) MyAppointmentsCommand {
	return MyAppointmentsCommand{
		logger:             logger,
		botGateway:         botGateway,
		tgUser:             tgUser,
		machine:            machine,
		userService:        userService,
		appointmentService: appointmentService,
	}
}

func (c MyAppointmentsCommand) Execute(ctx context.Context, messageDTO tg.MessageDTO) {
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser.TgID)
	appointments := c.appointmentService.GetAppointments(ctx, userEntity)

	if len(appointments) == 0 {
		c.botGateway.SendStartMessage(ctx, userEntity, messageDTO)
		c.machine.SetState(userEntity, state_machine.Start)
		return
	}

	c.botGateway.SendMyAppointmentsMessage(ctx, userEntity, appointments, messageDTO, 0)
	c.machine.SetState(userEntity, state_machine.ChooseAppointment)
}

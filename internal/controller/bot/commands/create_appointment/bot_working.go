package create_appointment

import (
	"context"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/pkg/client/telegram"
)

type CreateAppointmentCommand struct {
	logger             *slog.Logger
	bot                telegram.Bot
	botGateway         botGateway
	tgUser             tg.TgUserDTO
	userService        userService
	machine            *state_machine.UserStateMachine
	appointmentService appointmentService
	messageService     messageService
}

func NewCreateAppointmentCommand(
	logger *slog.Logger,
	bot telegram.Bot,
	botGateway botGateway,
	tgUser tg.TgUserDTO,
	userService userService,
	machine *state_machine.UserStateMachine,
	appointmentService appointmentService,
	messageService messageService,
) CreateAppointmentCommand {
	return CreateAppointmentCommand{
		logger:             logger,
		bot:                bot,
		botGateway:         botGateway,
		tgUser:             tgUser,
		userService:        userService,
		machine:            machine,
		appointmentService: appointmentService,
		messageService:     messageService,
	}
}

func (c CreateAppointmentCommand) Execute(ctx context.Context, messageDTO tg.MessageDTO) {
	// так как мы не изменяем бизнес сущность, а бот меняет состояние, то нахождение сущность в слое controllers некритично
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser.TgID)
	if userEntity.GetState() == nil {
		return
	}
	switch *userEntity.GetState() {
	case state_machine.Start:
		sentMessageId := c.botGateway.SendWaitMessage(ctx, userEntity, messageDTO, "wait speciality")

		specialities, err := c.appointmentService.GetSpecialities(ctx)
		if err != nil {
			return
		}

		translatedSpecialities, _, err := c.appointmentService.GetTranslatedSpecialities(ctx, userEntity, specialities, 0)
		if err != nil {
			return
		}

		c.botGateway.SendChooseSpecialityMessage(ctx, userEntity, messageDTO, sentMessageId, translatedSpecialities)
		go c.machine.SetState(userEntity, state_machine.ChooseSpeciality)
		go c.appointmentService.CreateDraftAppointment(ctx, userEntity.GetTgId())
	}
	return
}

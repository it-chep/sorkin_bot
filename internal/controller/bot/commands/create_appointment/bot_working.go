package create_appointment

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	case "":
		msgText, _ := c.messageService.GetMessage(ctx, userEntity, "wait speciality")
		sentMessageId := c.bot.SendMessageAndGetId(tgbotapi.NewMessage(c.tgUser.TgID, msgText), messageDTO)

		specialities, err := c.appointmentService.GetSpecialities(ctx)
		if err != nil {
			return
		}

		translatedSpecialities, _, err := c.appointmentService.GetTranslatedSpecialities(ctx, userEntity, specialities, 0)
		if err != nil {
			return
		}

		c.botGateway.SendChooseSpecialityMessage(ctx, sentMessageId, translatedSpecialities, userEntity, messageDTO)
		go c.machine.SetState(userEntity, state_machine.ChooseSpeciality)
		go c.appointmentService.CreateDraftAppointment(ctx, userEntity.GetTgId())
		return
	}
	return
}

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
	tgUser             tg.TgUserDTO
	userService        userService
	machine            *state_machine.UserStateMachine
	appointmentService appointmentService
	messageService     messageService
	botService         botService
}

func NewCreateAppointmentCommand(
	logger *slog.Logger,
	bot telegram.Bot,
	tgUser tg.TgUserDTO,
	userService userService,
	machine *state_machine.UserStateMachine,
	appointmentService appointmentService,
	messageService messageService,
	botService botService,
) CreateAppointmentCommand {
	return CreateAppointmentCommand{
		logger:             logger,
		bot:                bot,
		tgUser:             tgUser,
		userService:        userService,
		machine:            machine,
		appointmentService: appointmentService,
		messageService:     messageService,
		botService:         botService,
	}
}

func (c CreateAppointmentCommand) Execute(ctx context.Context, messageDTO tg.MessageDTO) {
	var msg tgbotapi.MessageConfig
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

		msgText, keyboard := c.botService.ConfigureGetSpecialityMessage(ctx, userEntity, translatedSpecialities, 0)
		msg = tgbotapi.NewMessage(c.tgUser.TgID, msgText)
		if keyboard.InlineKeyboard != nil {
			msg.ReplyMarkup = keyboard
		}
		c.bot.RemoveMessage(c.tgUser.TgID, sentMessageId)
		c.bot.SendMessage(msg, messageDTO)
		c.machine.SetState(userEntity, *userEntity.GetState(), state_machine.ChooseSpeciality)
		go c.appointmentService.CreateDraftAppointment(ctx, userEntity.GetTgId())
		return
	}

	return
}

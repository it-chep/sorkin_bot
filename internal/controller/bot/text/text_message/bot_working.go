package text_message

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/pkg/client/telegram"
)

type TextBotMessage struct {
	logger             *slog.Logger
	bot                telegram.Bot
	tgUser             tg.TgUserDTO
	machine            *state_machine.UserStateMachine
	userService        userService
	messageService     messageService
	appointmentService appointmentService
	botService         botService
}

func NewTextBotMessage(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO, machine *state_machine.UserStateMachine, userService userService, messageService messageService, appointmentService appointmentService, botService botService) TextBotMessage {
	return TextBotMessage{
		logger:             logger,
		bot:                bot,
		tgUser:             tgUser,
		machine:            machine,
		userService:        userService,
		messageService:     messageService,
		appointmentService: appointmentService,
		botService:         botService,
	}
}

func (c TextBotMessage) Execute(ctx context.Context, messageDTO tg.MessageDTO) {
	var _ tgbotapi.MessageConfig
	// так как мы не изменяем бизнес сущность, а бот меняет состояние, то нахождение сущность в слое controllers некритично
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser.TgID)

	switch *userEntity.GetState() {
	case "":
	case state_machine.GetPhone:
		if userEntity.GetPatientId() == nil {
			c.getPhone(ctx, userEntity, messageDTO)
		}
	case state_machine.GetName:
		if userEntity.GetPatientId() == nil {
			c.getName(ctx, userEntity, messageDTO)
		}
	case state_machine.GetBirthDate:
		if userEntity.GetPatientId() == nil {
			c.getBirthDate(ctx, userEntity, messageDTO)
		}
	}
}

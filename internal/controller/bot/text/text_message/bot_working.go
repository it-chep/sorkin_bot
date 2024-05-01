package text_message

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/bot/bot_interfaces"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/pkg/client/telegram"
)

type TextBotMessage struct {
	logger             *slog.Logger
	bot                telegram.Bot
	tgUser             tg.TgUserDTO
	machine            *state_machine.UserStateMachine
	userService        bot_interfaces.UserService
	messageService     bot_interfaces.MessageService
	appointmentService bot_interfaces.AppointmentService
	botService         bot_interfaces.BotService
}

func NewTextBotMessage(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO, machine *state_machine.UserStateMachine, userService bot_interfaces.UserService, messageService bot_interfaces.MessageService, appointmentService bot_interfaces.AppointmentService, botService bot_interfaces.BotService) TextBotMessage {
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
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser)
	switch userEntity.GetState() {
	case "":
	//	todo administrationHelp or start
	case state_machine.GetPhone:
		if userEntity.GetPatientId() != 0 {
			c.GetPhone(ctx, userEntity, messageDTO)
		}
	case state_machine.GetName:
		if userEntity.GetPatientId() != 0 {
			c.GetName(ctx, userEntity, messageDTO)
		}
	case state_machine.GetBirthDate:
		if userEntity.GetPatientId() != 0 {
			c.GetBirthDate(ctx, userEntity, messageDTO)
		}
	}

}

package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/pkg/client/telegram"
)

var languagesMap = map[string]bool{
	"RU": true,
	"EN": true,
	"PT": true,
}

const (
	ZeroOffset int = 0
)

type CallbackBotMessage struct {
	logger             *slog.Logger
	bot                telegram.Bot
	botGateway         botGateway
	tgUser             tg.TgUserDTO
	machine            *state_machine.UserStateMachine
	userService        userService
	messageService     messageService
	appointmentService appointmentService
}

func NewCallbackBot(
	logger *slog.Logger,
	bot telegram.Bot,
	botGateway botGateway,
	tgUser tg.TgUserDTO,
	machine *state_machine.UserStateMachine,
	userService userService,
	messageService messageService,
	appointmentService appointmentService,
) CallbackBotMessage {
	return CallbackBotMessage{
		logger:             logger,
		bot:                bot,
		botGateway:         botGateway,
		tgUser:             tgUser,
		machine:            machine,
		userService:        userService,
		messageService:     messageService,
		appointmentService: appointmentService,
	}
}

// Execute место связи telegram и бизнес логи
func (c *CallbackBotMessage) Execute(ctx context.Context, messageDTO tg.MessageDTO, callbackData string) {

	userEntity, _ := c.userService.GetUser(ctx, c.tgUser.TgID)

	switch *userEntity.GetState() {
	case state_machine.Start:
		c.mainMenu(ctx, messageDTO, userEntity, callbackData)
	case state_machine.ChooseSpeciality:
		c.chooseSpeciality(ctx, messageDTO, userEntity, callbackData)
	case state_machine.ChooseAppointment:
		c.getAppointmentDetail(ctx, messageDTO, callbackData)
	case state_machine.ChooseSchedule:
		c.chooseSchedules(ctx, messageDTO, userEntity, callbackData)
	case state_machine.CreateAppointment:
		c.preCreateAppointment(ctx, messageDTO, userEntity, callbackData)
	case state_machine.FastAppointment:
		c.fastAppointment(ctx, messageDTO, userEntity, callbackData)
	case state_machine.ChooseDoctor:
		c.chooseDoctor(ctx, messageDTO, userEntity, callbackData)
	case state_machine.DetailMyAppointment:
		c.detailMyAppointment(ctx, messageDTO, callbackData)
	case state_machine.ChooseLanguage:
		c.chooseLanguage(ctx, messageDTO, callbackData)
	}
}

func (c *CallbackBotMessage) chooseLanguage(ctx context.Context, messageDTO tg.MessageDTO, callbackData string) {
	userEntity, err := c.userService.ChangeLanguage(ctx, c.tgUser, callbackData)
	if err != nil {
		return
	}

	userEntity.SetLanguageCode(callbackData)
	msgText, err := c.messageService.GetMessage(ctx, userEntity, "successfully changed language")

	c.bot.RemoveMessage(userEntity.GetTgId(), int(messageDTO.MessageID))

	if err != nil {
		c.logger.Error(fmt.Sprintf("error: %s,  place: CallbackBotMessage/chooseLanguage", err))
		c.botGateway.SendError(ctx, userEntity, messageDTO)
		return
	}

	c.bot.SendMessage(
		tgbotapi.NewMessage(c.tgUser.TgID, fmt.Sprintf(msgText, callbackData)),
		messageDTO,
	)

	c.botGateway.SendStartMessage(ctx, userEntity, messageDTO)
	go c.machine.SetState(userEntity, state_machine.Start)
}

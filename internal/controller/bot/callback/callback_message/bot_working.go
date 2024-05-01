package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/bot"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/message"
	"sorkin_bot/pkg/client/telegram"
)

var languagesMap = map[string]bool{
	"RU": true,
	"EN": true,
	"PT": true,
}

type CallbackBotMessage struct {
	logger             *slog.Logger
	bot                telegram.Bot
	tgUser             tg.TgUserDTO
	machine            *state_machine.UserStateMachine
	userService        bot.UserService
	messageService     bot.MessageService
	appointmentService bot.AppointmentService
}

func NewCallbackBot(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO, machine *state_machine.UserStateMachine, userService bot.UserService, messageService bot.MessageService, appointmentService bot.AppointmentService) CallbackBotMessage {
	return CallbackBotMessage{
		logger:             logger,
		bot:                bot,
		tgUser:             tgUser,
		machine:            machine,
		userService:        userService,
		messageService:     messageService,
		appointmentService: appointmentService,
	}
}

// Execute место связи telegram и бизнес логи
func (c *CallbackBotMessage) Execute(ctx context.Context, messageDTO tg.MessageDTO, callbackData string) {
	op := "sorkin_bot.internal.controller.bot.callback.callback_message.bot_working.Execute"
	var msg tgbotapi.MessageConfig
	var msgText string
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser)
	if _, ok := languagesMap[callbackData]; ok {
		userEntity, err := c.userService.ChangeLanguage(ctx, c.tgUser, callbackData)
		if err != nil {
			return
		}
		userEntity.SetLanguageCode(callbackData)
		msgText, err = c.messageService.GetMessage(ctx, userEntity, "successfully changed language")
		if err != nil {
			c.logger.Error(fmt.Sprintf("error: %s,  place: %s", err, op))
			msgText = message.ServerError
		}
	}

	switch userEntity.GetState() {
	case state_machine.ChooseAppointment:
		c.GetAppointmentDetail(ctx, messageDTO, callbackData)
	case state_machine.ChooseSchedule:
		fmt.Println(1212)
	case state_machine.ChooseDoctor:
		fmt.Println(232323)
	}

	msg = tgbotapi.NewMessage(c.tgUser.TgID, msgText)

	_, err := c.bot.Bot.Send(msg)
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s", err))
	}
}

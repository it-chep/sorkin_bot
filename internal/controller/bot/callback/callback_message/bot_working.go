package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/message"
	"sorkin_bot/pkg/client/telegram"
	"strconv"
	"strings"
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
	userService        UserService
	messageService     MessageService
	appointmentService AppointmentService
	botService         BotService
}

func NewCallbackBot(
	logger *slog.Logger,
	bot telegram.Bot,
	tgUser tg.TgUserDTO,
	machine *state_machine.UserStateMachine,
	userService UserService,
	messageService MessageService,
	appointmentService AppointmentService,
	botService BotService,
) CallbackBotMessage {
	return CallbackBotMessage{
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

// Execute место связи telegram и бизнес логи
func (c *CallbackBotMessage) Execute(ctx context.Context, messageDTO tg.MessageDTO, callbackData string) {
	op := "sorkin_bot.internal.controller.bot.callback.callback_message.bot_working.Execute"
	var msg tgbotapi.MessageConfig
	var msgText string
	var err error

	userEntity, _ := c.userService.GetUser(ctx, c.tgUser)
	if _, ok := languagesMap[callbackData]; ok {
		userEntity, err = c.userService.ChangeLanguage(ctx, c.tgUser, callbackData)
		if err != nil {
			return
		}
		userEntity.SetLanguageCode(callbackData)
		msgText, err = c.messageService.GetMessage(ctx, userEntity, "successfully changed language")
		if err != nil {
			c.logger.Error(fmt.Sprintf("error: %s,  place: %s", err, op))
			msgText = message.ServerError
		}
		msg = tgbotapi.NewMessage(c.tgUser.TgID, msgText)
		c.bot.SendMessage(msg, messageDTO)
	}
	c.logger.Info(fmt.Sprintf("ENTER TO Execute SPEC callback %s", callbackData))

	switch userEntity.GetState() {
	case state_machine.ChooseSpeciality:
		c.chooseSpeciality(ctx, messageDTO, userEntity, callbackData)
	case state_machine.ChooseAppointment:
		c.GetAppointmentDetail(ctx, messageDTO, callbackData)
	case state_machine.ChooseSchedule:
		c.GetSchedules(ctx, messageDTO, callbackData)
	case state_machine.DetailMyAppointment:
		c.DetailMyAppointment(ctx, messageDTO, callbackData)
	}
}

func (c *CallbackBotMessage) chooseSpeciality(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {

	if strings.Contains(callbackData, "offset") {
		c.moreLessSpeciality(ctx, messageDTO, userEntity, callbackData)
	} else {
		c.GetDoctors(ctx, messageDTO, callbackData)
	}
}

func (c *CallbackBotMessage) moreLessSpeciality(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	specialities, err := c.appointmentService.GetSpecialities(ctx)
	if err != nil {
		return
	}
	offset, _ := strconv.Atoi(strings.Split(callbackData, "_")[1])
	if strings.Contains(callbackData, ">") {
		offset += 10
	} else {
		offset -= 10
	}
	translatedSpecialities, _, err := c.appointmentService.GetTranslatedSpecialities(ctx, userEntity, specialities, offset)
	if err != nil {
		return
	}
	msgText, keyboard := c.botService.ConfigureGetSpecialityMessage(ctx, userEntity, translatedSpecialities, offset)
	msg := tgbotapi.NewMessage(c.tgUser.TgID, msgText)
	msg.ReplyMarkup = keyboard

	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
	c.bot.SendMessage(msg, messageDTO)
}

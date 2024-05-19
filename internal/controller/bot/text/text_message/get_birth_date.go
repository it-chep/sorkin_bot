package text_message

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/message"
	"strconv"
	"strings"
	"time"
)

func (c TextBotMessage) getBirthDate(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	var msg tgbotapi.MessageConfig
	var err error
	if c.validateBirthDateMessage(messageDTO.Text) {
		user, err = c.userService.UpdateBirthDate(ctx, c.tgUser, messageDTO.Text)
		if err != nil {
			msg = tgbotapi.NewMessage(c.tgUser.TgID, message.ServerError)
			c.bot.SendMessage(msg, messageDTO)
			return
		}
	} else {
		messageText, _ := c.messageService.GetMessage(ctx, user, "invalid birth date")
		msg = tgbotapi.NewMessage(c.tgUser.TgID, messageText)
		c.bot.SendMessage(msg, messageDTO)
		return
	}

	messageText, _ := c.messageService.GetMessage(ctx, user, "ready to appointment")
	msg = tgbotapi.NewMessage(c.tgUser.TgID, messageText)
	c.bot.SendMessage(msg, messageDTO)
	c.machine.SetState(user, state_machine.GetBirthDate, state_machine.CreateAppointment)

	if c.appointmentService.GetPatient(ctx, user) {
	} else {
		c.appointmentService.CreatePatient(ctx, user)
	}
}

func (c TextBotMessage) validateBirthDateMessage(birthDate string) (valid bool) {
	currentTime := time.Now()
	dateItems := strings.Split(birthDate, ".")
	validDateToday := false

	validLength := len(birthDate) == 10

	validDateItemsLength := len(dateItems) == 3

	if validDateItemsLength {
		intDay, err := strconv.Atoi(dateItems[0])
		if err != nil {
			return false
		}
		intMonth, err := strconv.Atoi(dateItems[1])
		if err != nil {
			return false
		}
		intYear, err := strconv.Atoi(dateItems[2])
		if err != nil {
			return false
		}

		unvalidatedDate := time.Date(intYear, time.Month(intMonth), intDay, 0, 0, 0, 0, time.UTC)
		validDateToday = unvalidatedDate.Before(currentTime)
	}
	if validLength && validDateItemsLength && validDateToday {
		return true
	}

	return false
}

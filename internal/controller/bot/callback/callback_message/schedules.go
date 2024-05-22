package callback

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"strconv"
	"strings"
)

func (c *CallbackBotMessage) getSchedules(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	var msgText string

	msgText, err := c.messageService.GetMessage(ctx, userEntity, "load doctor schedules")
	msg := tgbotapi.NewMessage(c.tgUser.TgID, msgText)
	messageId := c.bot.SendMessageAndGetId(msg, messageDTO)

	doctorId, err := strconv.Atoi(callbackData)
	if err != nil {
		return
	}

	schedules, err := c.appointmentService.GetSchedules(ctx, userEntity, &doctorId)
	if err != nil {
		return
	}

	msgText, keyboard := c.botService.ConfigureGetScheduleMessage(ctx, userEntity, schedules, 0)

	c.bot.RemoveMessage(c.tgUser.TgID, messageId)
	msg = tgbotapi.NewMessage(c.tgUser.TgID, msgText)
	msg.ReplyMarkup = keyboard
	c.bot.SendMessage(msg, messageDTO)
}

func (c *CallbackBotMessage) chooseSchedules(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	if strings.Contains(callbackData, "offset") {
		c.moreLessSchedules(ctx, messageDTO, userEntity, callbackData)
	} else if strings.Contains(callbackData, "schedule") {
		c.saveDraftAppointment(ctx, messageDTO, userEntity, callbackData)
	} else {
		c.getSchedules(ctx, messageDTO, userEntity, callbackData)
	}
}

func (c *CallbackBotMessage) moreLessSchedules(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	offset, _ := strconv.Atoi(strings.Split(callbackData, "_")[1])
	if strings.Contains(callbackData, ">") {
		offset += 10
	} else {
		offset -= 10
	}

	schedules, err := c.appointmentService.GetSchedules(ctx, userEntity, nil)
	if err != nil {
		return
	}
	msgText, keyboard := c.botService.ConfigureGetScheduleMessage(ctx, userEntity, schedules, offset)
	msg := tgbotapi.NewMessage(c.tgUser.TgID, msgText)
	// todo протестить не будет ли бага
	msg.ReplyMarkup = keyboard

	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
	c.bot.SendMessage(msg, messageDTO)
}

func (c *CallbackBotMessage) saveDraftAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	scheduleItems := strings.Split(callbackData, "_")
	doctorId, _ := strconv.Atoi(scheduleItems[1])
	fullTimeStart := scheduleItems[2]
	fullTimeEnd := scheduleItems[3]
	date := scheduleItems[4]
	c.appointmentService.UpdateDraftAppointmentDate(ctx, userEntity.GetTgId(), fullTimeStart, fullTimeEnd, date)

	if userEntity.GetPhone() == nil {
		c.machine.SetState(userEntity, *userEntity.GetState(), state_machine.GetPhone)

		msgText, keyboard := c.botService.ConfigureGetPhoneMessage(ctx, userEntity)
		msg := tgbotapi.NewMessage(c.tgUser.TgID, msgText)
		msg.ReplyMarkup = keyboard
		c.bot.SendMessage(msg, messageDTO)
		return
	}

	msgText, keyboard := c.botService.ConfigureConfirmAppointmentMessage(ctx, userEntity, doctorId)
	msg := tgbotapi.NewMessage(c.tgUser.TgID, msgText)
	msg.ReplyMarkup = keyboard
	c.bot.SendMessage(msg, messageDTO)

	c.machine.SetState(userEntity, *userEntity.GetState(), state_machine.CreateAppointment)
}

package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"strconv"
	"strings"
)

func (c *CallbackBotMessage) getSchedules(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	var msgText string
	var err error
	var msg tgbotapi.MessageConfig

	callbackDataItems := strings.Split(callbackData, "_")
	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))

	msgText, err = c.messageService.GetMessage(ctx, userEntity, "your doctor")
	msg = tgbotapi.NewMessage(c.tgUser.TgID, fmt.Sprintf(msgText, callbackDataItems[1]))
	c.bot.SendMessage(msg, messageDTO)

	sentMessageId := c.botGateway.SendWaitMessage(ctx, userEntity, messageDTO, "wait schedules")

	doctorId, err := strconv.Atoi(callbackDataItems[0])
	if err != nil {
		return
	}

	schedules, err := c.appointmentService.GetSchedules(ctx, userEntity, &doctorId)

	if err != nil {
		msgText, err = c.messageService.GetMessage(ctx, userEntity, "empty schedules")
		msg = tgbotapi.NewMessage(c.tgUser.TgID, msgText)
		c.bot.SendMessage(msg, messageDTO)
		c.machine.SetState(userEntity, state_machine.ChooseDoctor)

		draftAppointmentEntity, _ := c.appointmentService.GetDraftAppointment(ctx, userEntity.GetTgId())

		c.getDoctors(ctx, messageDTO, userEntity, *draftAppointmentEntity.GetSpecialityId())
		return
	}

	c.bot.RemoveMessage(c.tgUser.TgID, sentMessageId)
	c.botGateway.SendSchedulesMessage(ctx, userEntity, messageDTO, schedules, ZeroOffset)
	c.machine.SetState(userEntity, state_machine.ChooseSchedule)
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

	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
	c.botGateway.SendSchedulesMessage(ctx, userEntity, messageDTO, schedules, offset)
}

func (c *CallbackBotMessage) saveDraftAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	scheduleItems := strings.Split(callbackData, "_")
	doctorId, _ := strconv.Atoi(scheduleItems[1])
	fullTimeStart := scheduleItems[2]
	fullTimeEnd := scheduleItems[3]
	date := scheduleItems[4]

	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))

	c.appointmentService.UpdateDraftAppointmentDate(ctx, userEntity.GetTgId(), fullTimeStart, fullTimeEnd, date)

	if userEntity.GetPhone() == nil {
		c.botGateway.SendGetPhoneMessage(ctx, userEntity, messageDTO)
		c.machine.SetState(userEntity, state_machine.GetPhone)
		return
	}

	c.botGateway.SendConfirmAppointmentMessage(ctx, userEntity, messageDTO, doctorId)
	c.machine.SetState(userEntity, state_machine.CreateAppointment)
}

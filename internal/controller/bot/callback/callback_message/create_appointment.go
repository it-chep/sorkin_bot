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

func (c *CallbackBotMessage) preCreateAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	if callbackData == "confirm_appointment" {
		c.confirmAppointment(ctx, messageDTO, userEntity)
	} else if callbackData == "reject_appointment" {
		c.appointmentService.CleanDraftAppointment(ctx, userEntity.GetTgId())
		c.rejectAppointment(ctx, messageDTO, userEntity)
	} else if strings.Contains(callbackData, "doc_info") {
		c.getDoctorInfo(ctx, messageDTO, userEntity, callbackData)
	}
}

func (c *CallbackBotMessage) rejectAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User) {
	msgText, err := c.messageService.GetMessage(ctx, userEntity, "have rejected appointment")
	if err != nil {
		return
	}
	msg := tgbotapi.NewMessage(userEntity.GetTgId(), msgText)
	c.bot.SendMessage(msg, messageDTO)
	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
	c.botGateway.SendStartMessage(ctx, userEntity, messageDTO)
	c.machine.SetState(userEntity, state_machine.Start)
}

func (c *CallbackBotMessage) confirmAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User) {
	draftAppointmentEntity, err := c.appointmentService.GetDraftAppointment(ctx, userEntity.GetTgId())
	if err != nil {
		return
	}
	if draftAppointmentEntity.GetDoctorId() == nil {
		return
	}

	//todo херня, надо улучшить
	appointmentString := fmt.Sprintf("doctorId_%d__timeStart_%s__timeEnd_%s",
		*draftAppointmentEntity.GetDoctorId(),
		*draftAppointmentEntity.GetTimeStart(),
		*draftAppointmentEntity.GetTimeEnd(),
	)

	appointmentId := c.appointmentService.CreateAppointment(ctx, userEntity, appointmentString)
	if appointmentId != nil {
		c.appointmentService.UpdateDraftAppointmentStatus(ctx, userEntity.GetTgId(), *appointmentId)
	}

	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))

	msgText, _ := c.messageService.GetMessage(ctx, userEntity, "successfully created appointment")
	msg := tgbotapi.NewMessage(c.tgUser.TgID, fmt.Sprintf(msgText, *draftAppointmentEntity.GetTimeStart()))
	c.bot.SendMessage(msg, messageDTO)

	c.botGateway.SendStartMessage(ctx, userEntity, messageDTO)
	c.machine.SetState(userEntity, state_machine.Start)
}

func (c *CallbackBotMessage) fastAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	if strings.Contains(callbackData, "fast_") {
		c.bot.RemoveMessage(userEntity.GetTgId(), int(messageDTO.MessageID))

		items := strings.Split(callbackData, "__")
		specialityId, err := strconv.Atoi(items[1])
		if err != nil {
			return
		}
		doctorId, err := strconv.Atoi(items[2])
		if err != nil {
			return
		}
		timeStart := items[3]
		timeEnd := items[4]
		c.appointmentService.FastUpdateDraftAppointment(ctx, userEntity.GetTgId(), specialityId, doctorId, timeStart, timeEnd)

		if userEntity.GetPhone() != nil {
			c.botGateway.SendConfirmAppointmentMessage(ctx, userEntity, messageDTO, doctorId)
			c.machine.SetState(userEntity, state_machine.CreateAppointment)
		} else {
			c.botGateway.SendGetPhoneMessage(ctx, userEntity, messageDTO)
			c.machine.SetState(userEntity, state_machine.GetPhone)
		}
	}
}

func (c *CallbackBotMessage) getDoctorInfo(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	callbackItems := strings.Split(callbackData, "_")
	doctorId, _ := strconv.Atoi(callbackItems[2])
	c.botGateway.SendDoctorInfoMessage(ctx, userEntity, messageDTO, int(messageDTO.MessageID), doctorId)
	c.machine.SetState(userEntity, state_machine.GetDoctorInfo)
}

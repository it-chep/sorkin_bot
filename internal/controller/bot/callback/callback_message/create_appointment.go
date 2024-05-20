package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"strings"
)

func (c *CallbackBotMessage) preCreateAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	if callbackData == "confirm_appointment" {
		c.confirmAppointment(ctx, messageDTO, userEntity)
	} else if callbackData == "reject_appointment" {
		c.appointmentService.CleanDraftAppointment(ctx, userEntity.GetTgId())
		c.rejectAppointment(ctx, messageDTO, userEntity)
	}
}

func (c *CallbackBotMessage) rejectAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User) {
	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
	msgText, _ := c.messageService.GetMessage(ctx, userEntity, "Start")
	msg := tgbotapi.NewMessage(c.tgUser.TgID, msgText)
	c.bot.SendMessage(msg, messageDTO)
	c.machine.SetState(userEntity, *userEntity.GetState(), state_machine.Start)
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
	c.machine.SetState(userEntity, *userEntity.GetState(), state_machine.Start)

	msgText, _ := c.messageService.GetMessage(ctx, userEntity, "successfully created appointment")
	msg := tgbotapi.NewMessage(c.tgUser.TgID, fmt.Sprintf(msgText, *draftAppointmentEntity.GetTimeStart()))
	c.bot.SendMessage(msg, messageDTO)

	msgText, _ = c.messageService.GetMessage(ctx, userEntity, "Start")
	msg = tgbotapi.NewMessage(c.tgUser.TgID, msgText)
	c.bot.SendMessage(msg, messageDTO)
	c.machine.SetState(userEntity, *userEntity.GetState(), state_machine.Start)
}

func (c *CallbackBotMessage) fastAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	if strings.Contains(callbackData, "fast_") {
		//items := strings.Split(callbackData, "__")
		//doctorId, _ := strconv.Atoi(items[1])
		//timeStart := items[2]
		//timeEnd := items[3]
	}
}

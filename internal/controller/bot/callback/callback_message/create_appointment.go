package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/entity/user/state_machine"
)

func (c *CallbackBotMessage) preCreateAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	if callbackData == "confirm_appointment" {
		c.appointmentService.UpdateDraftAppointmentStatus(ctx, userEntity.GetTgId())
		c.confirmAppointment(ctx, messageDTO, userEntity)
	} else if callbackData == "reject_appointment" {
		c.appointmentService.CleanDraftAppointment(ctx, userEntity.GetTgId())
		c.rejectAppointment(ctx, messageDTO, userEntity)
	}
}

func (c *CallbackBotMessage) rejectAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User) {
	msgText, err := c.messageService.GetMessage(ctx, userEntity, "Start")
	if err != nil {
		return
	}
	msg := tgbotapi.NewMessage(c.tgUser.TgID, msgText)
	c.bot.SendMessage(msg, messageDTO)
	c.machine.SetState(userEntity, *userEntity.GetState(), state_machine.Start)
}

func (c *CallbackBotMessage) confirmAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User) {
	draftAppointmentEntity, err := c.appointmentService.GetDraftAppointment(ctx, userEntity.GetTgId())
	if err != nil {
		return
	}
	//todo херня, надо улучшить
	appointmentString := fmt.Sprintf("doctorId_%d__timeStart_%s__timeEnd_%s",
		draftAppointmentEntity.GetDoctorId(),
		*draftAppointmentEntity.GetTimeStart(),
		*draftAppointmentEntity.GetTimeEnd(),
	)

	c.appointmentService.CreateAppointment(ctx, userEntity, appointmentString)
}

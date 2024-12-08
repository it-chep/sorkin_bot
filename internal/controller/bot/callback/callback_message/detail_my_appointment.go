package callback

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"strconv"
	"strings"
)

func (c *CallbackBotMessage) getAppointmentDetail(ctx context.Context, messageDTO tg.MessageDTO, callbackData string) {
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser.TgID)

	if strings.Contains(callbackData, ">") || strings.Contains(callbackData, "<") {
		c.moreLessMyAppointments(ctx, messageDTO, userEntity, callbackData)
		return
	}

	c.bot.RemoveMessage(userEntity.GetTgId(), int(messageDTO.MessageID))

	appointmentId, err := strconv.Atoi(strings.Split(callbackData, "_")[1])
	if err != nil {
		c.botGateway.SendError(ctx, userEntity, messageDTO)
		return
	}

	appointmentEntity := c.appointmentService.GetAppointmentDetail(ctx, userEntity, appointmentId)
	if appointmentEntity.GetAppointmentId() != 0 {
		c.botGateway.SendDetailAppointmentMessage(ctx, userEntity, messageDTO, appointmentEntity)
	} else {
		c.botGateway.SendEmptyAppointments(ctx, userEntity, messageDTO)
	}

	c.machine.SetState(userEntity, state_machine.DetailMyAppointment)
}

func (c *CallbackBotMessage) detailMyAppointment(ctx context.Context, messageDTO tg.MessageDTO, callbackData string) {
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser.TgID)
	dataItems := strings.Split(callbackData, "_")

	if strings.Contains(callbackData, "doc_info") {
		c.getDoctorInfo(ctx, messageDTO, userEntity, callbackData)
		return
	}

	if len(dataItems) == 1 {
		c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
		c.botGateway.SendStartMessage(ctx, userEntity, messageDTO)
		c.machine.SetState(userEntity, state_machine.Start)
		return
	}

	appointmentId, err := strconv.Atoi(dataItems[1])
	if err != nil {
		return
	}
	if dataItems[0] == "cancel" {
		c.cancelAppointment(ctx, messageDTO, userEntity, appointmentId)
	} else if dataItems[0] == "reschedule" {
		c.rescheduleAppointment(ctx, messageDTO, userEntity, appointmentId)
	}
}

func (c *CallbackBotMessage) cancelAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, appointmentId int) {
	if c.appointmentService.CancelAppointment(ctx, userEntity, appointmentId) {
		c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
		c.botGateway.SendStartMessage(ctx, userEntity, messageDTO)
		c.machine.SetState(userEntity, state_machine.Start)
	}
}

func (c *CallbackBotMessage) rescheduleAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, appointmentId int) {
	c.appointmentService.CancelAppointment(ctx, userEntity, appointmentId)
	draftAppointmentEntity, err := c.appointmentService.GetDraftAppointmentByAppointmentId(ctx, appointmentId)
	if err != nil {
		return
	}
	doctorsMap := c.appointmentService.GetDoctorsBySpecialityId(ctx, userEntity.GetTgId(), ZeroOffset, draftAppointmentEntity.GetSpecialityId())
	c.appointmentService.FastUpdateDraftAppointment(
		ctx, userEntity.GetTgId(),
		*draftAppointmentEntity.GetSpecialityId(),
		*draftAppointmentEntity.GetDoctorId(),
		*draftAppointmentEntity.GetTimeStart(),
		*draftAppointmentEntity.GetTimeEnd(),
	)
	c.botGateway.SendGetDoctorsMessage(ctx, userEntity, messageDTO, doctorsMap, ZeroOffset)
	c.machine.SetState(userEntity, state_machine.ChooseDoctor)
}

func (c *CallbackBotMessage) moreLessMyAppointments(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	offset, _ := strconv.Atoi(strings.Split(callbackData, "_")[1])
	if strings.Contains(callbackData, ">") {
		offset += 10
	} else {
		offset -= 10
	}

	appointments := c.appointmentService.GetAppointments(ctx, userEntity)
	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))

	c.botGateway.SendMyAppointmentsMessage(ctx, userEntity, appointments, messageDTO, offset)
}

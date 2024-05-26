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

	c.bot.RemoveMessage(userEntity.GetTgId(), int(messageDTO.MessageID))

	appointmentId, err := strconv.Atoi(strings.Split(callbackData, "_")[1])
	if err != nil {
		c.botGateway.SendError(ctx, userEntity, messageDTO)
		return
	}

	appointmentEntity := c.appointmentService.GetAppointmentDetail(ctx, userEntity, appointmentId)
	c.logger.Info("appointmentDetail and callbackData", appointmentEntity, callbackData, appointmentId)
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
	if dataItems[0] == "cancel" {
		appointmentId, _ := strconv.Atoi(dataItems[1])
		c.cancelAppointment(ctx, messageDTO, userEntity, appointmentId)
		//} else if dataItems[0] == "reschedule" {
		//	appointmentId, _ := strconv.Atoi(dataItems[1])
		//	c.rescheduleAppointment(ctx, messageDTO, userEntity, "", appointmentId)
	} else if strings.Contains(callbackData, "doc_info") {
		c.getDoctorInfo(ctx, messageDTO, userEntity, callbackData)
	} else if dataItems[0] == "exit" {
		c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
		c.botGateway.SendStartMessage(ctx, userEntity, messageDTO)
		c.machine.SetState(userEntity, state_machine.Start)
	}
}

func (c *CallbackBotMessage) cancelAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, appointmentId int) {
	if c.appointmentService.CancelAppointment(ctx, userEntity, appointmentId) {
		c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
		c.botGateway.SendStartMessage(ctx, userEntity, messageDTO)
		c.machine.SetState(userEntity, state_machine.Start)
	}
}

//func (c *CallbackBotMessage) rescheduleAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, movedTo string, appointmentId int) {
//	//c.appointmentService.RescheduleAppointment(ctx, appointmentId, movedTo)
//}
//
//func (c *CallbackBotMessage) moveAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, movedTo string, appointmentId int) {}

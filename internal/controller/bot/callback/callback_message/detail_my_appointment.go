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

func (c *CallbackBotMessage) detailMyAppointment(ctx context.Context, messageDTO tg.MessageDTO, callbackData string) {
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser.TgID)
	dataItems := strings.Split(callbackData, "_")
	if dataItems[0] == "cancel" {
		appointmentId, _ := strconv.Atoi(dataItems[1])
		c.cancelAppointment(ctx, messageDTO, userEntity, appointmentId)
		//} else if dataItems[0] == "reschedule" {
		//	appointmentId, _ := strconv.Atoi(dataItems[1])
		//	c.rescheduleAppointment(ctx, messageDTO, userEntity, "", appointmentId)
	} else {
		return
	}
}

func (c *CallbackBotMessage) cancelAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, appointmentId int) {
	if c.appointmentService.CancelAppointment(ctx, userEntity, appointmentId) {
		c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
		// todo докрутить сообщения
		msgText, _ := c.messageService.GetMessage(ctx, userEntity, "")
		c.bot.SendMessage(tgbotapi.NewMessage(c.tgUser.TgID, msgText), messageDTO)
		c.machine.SetState(userEntity, *userEntity.GetState(), state_machine.Start)
	} else {

	}

}

//func (c *CallbackBotMessage) rescheduleAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, movedTo string, appointmentId int) {
//	//c.appointmentService.RescheduleAppointment(ctx, appointmentId, movedTo)
//}
//
//func (c *CallbackBotMessage) moveAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, movedTo string, appointmentId int) {}

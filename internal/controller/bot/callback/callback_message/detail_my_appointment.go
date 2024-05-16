package callback

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
	"strconv"
	"strings"
)

func (c *CallbackBotMessage) DetailMyAppointment(ctx context.Context, messageDTO tg.MessageDTO, callbackData string) {
	dataItems := strings.Split(callbackData, "_")
	if dataItems[0] == "cancel" {
		appointmentId, _ := strconv.Atoi(dataItems[1])
		c.cancelAppointment(ctx, appointmentId)
	} else if dataItems[0] == "reschedule" {
		appointmentId, _ := strconv.Atoi(dataItems[1])
		c.rescheduleAppointment(ctx, "", appointmentId)
	} else {
		return
	}
}

func (c *CallbackBotMessage) cancelAppointment(ctx context.Context, appointmentId int) {
	//c.appointmentService.CancelAppointment(ctx, appointmentId)
}

func (c *CallbackBotMessage) rescheduleAppointment(ctx context.Context, movedTo string, appointmentId int) {
	//c.appointmentService.RescheduleAppointment(ctx, appointmentId, movedTo)
}

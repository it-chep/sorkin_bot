package callback

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
	"strconv"
)

func (c *CallbackBotMessage) GetSchedules(ctx context.Context, messageDTO tg.MessageDTO, callbackData string) {
	doctorId, err := strconv.Atoi(callbackData)
	if err != nil {
		return
	}

	c.appointmentService.GetSchedules(ctx, doctorId)

}

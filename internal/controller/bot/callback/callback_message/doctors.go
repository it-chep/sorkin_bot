package callback

import (
	"context"
	"fmt"
	"sorkin_bot/internal/controller/dto/tg"
	"strconv"
)

func (c *CallbackBotMessage) GetDoctors(ctx context.Context, messageDTO tg.MessageDTO, callbackData string) {
	specialityId, err := strconv.Atoi(callbackData)
	if err != nil {
		return
	}

	doctors := c.appointmentService.GetDoctors(ctx, specialityId)
	c.logger.Info(fmt.Sprintf("%s", doctors))

}

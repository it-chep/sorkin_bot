package callback

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"strconv"
	"strings"
)

func (c *CallbackBotMessage) chooseSpeciality(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {

	if strings.Contains(callbackData, "offset") {
		c.moreLessSpeciality(ctx, messageDTO, userEntity, callbackData)
	} else {
		specialityId, _ := strconv.Atoi(callbackData)
		c.getDoctors(ctx, messageDTO, userEntity, specialityId)
		go c.appointmentService.UpdateDraftAppointmentIntField(ctx, userEntity.GetTgId(), specialityId, "speciality_id")
	}

}

func (c *CallbackBotMessage) moreLessSpeciality(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	specialities, err := c.appointmentService.GetSpecialities(ctx)
	if err != nil {
		return
	}

	offset := 0
	if callbackData != "" {
		offset, _ = strconv.Atoi(strings.Split(callbackData, "_")[1])
		if strings.Contains(callbackData, ">") {
			offset += 10
		} else {
			offset -= 10
		}
	}

	translatedSpecialities, _, err := c.appointmentService.GetTranslatedSpecialities(ctx, userEntity, specialities, offset)
	if err != nil {
		return
	}

	if callbackData != "" {
		c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
	}

	c.botGateway.SendSpecialityMessage(ctx, userEntity, messageDTO, translatedSpecialities, offset)
}

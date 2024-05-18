package callback

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"strconv"
	"strings"
)

func (c *CallbackBotMessage) GetDoctors(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, specialityId int) {
	doctors := c.appointmentService.GetDoctors(ctx, userEntity.GetTgId(), 0, &specialityId)

	msgText, keyboard := c.botService.ConfigureGetDoctorMessage(ctx, userEntity, doctors, 0)
	msg := tgbotapi.NewMessage(c.tgUser.TgID, msgText)

	if keyboard.InlineKeyboard != nil {
		msg.ReplyMarkup = keyboard
		c.bot.SendMessage(msg, messageDTO)
	} else {
		c.bot.SendMessage(msg, messageDTO)
		c.moreLessSpeciality(ctx, messageDTO, userEntity, "")
	}
}

func (c *CallbackBotMessage) chooseDoctor(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	if strings.Contains(callbackData, "offset") {
		c.moreLessDoctors(ctx, messageDTO, userEntity, callbackData)
	} else {
		c.GetSchedules(ctx, messageDTO, callbackData)
	}
}

func (c *CallbackBotMessage) moreLessDoctors(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	offset, _ := strconv.Atoi(strings.Split(callbackData, "_")[1])
	if strings.Contains(callbackData, ">") {
		offset += 10
	} else {
		offset -= 10
	}

	doctors := c.appointmentService.GetDoctors(ctx, userEntity.GetTgId(), offset, nil)

	msgText, keyboard := c.botService.ConfigureGetDoctorMessage(ctx, userEntity, doctors, offset)
	msg := tgbotapi.NewMessage(c.tgUser.TgID, msgText)
	// todo протестить не будет ли бага
	msg.ReplyMarkup = keyboard

	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
	c.bot.SendMessage(msg, messageDTO)
}

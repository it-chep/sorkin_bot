package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/message"
	"strconv"
)

func (c *CallbackBotMessage) GetAppointmentDetail(ctx context.Context, messageDTO tg.MessageDTO, callbackData string) {
	var msg tgbotapi.MessageConfig
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser)
	//todo возможно добавить сообщение, что я загружаю ваши записи, пожалуйста подождите
	//msg = tgbotapi.NewMessage(c.tgUser.TgID)
	//c.bot.SendMessage(msg, messageDTO)

	appointmentId, err := strconv.Atoi(callbackData)
	if err != nil {
		msg = tgbotapi.NewMessage(c.tgUser.TgID, message.ServerError)
		c.bot.SendMessage(msg, messageDTO)
		return
	}

	appointmentEntity := c.appointmentService.GetAppointmentDetail(ctx, userEntity, appointmentId)

	if appointmentEntity.GetAppointmentId() != 0 {
		msg = tgbotapi.NewMessage(c.tgUser.TgID, fmt.Sprintf("Запись на %s %d", appointmentEntity.GetTimeStart(), appointmentEntity.GetAppointmentId()))
		cancelText, err := c.messageService.GetMessage(ctx, userEntity, "cancel appointment button")
		if err != nil {
			msg = tgbotapi.NewMessage(c.tgUser.TgID, message.ServerError)
			c.bot.SendMessage(msg, messageDTO)
			return
		}
		rescheduleText, err := c.messageService.GetMessage(ctx, userEntity, "reschedule appointment button")
		if err != nil {
			msg = tgbotapi.NewMessage(c.tgUser.TgID, message.ServerError)
			c.bot.SendMessage(msg, messageDTO)
			return
		}
		// формируем клавиатуру действий с онлайн записью
		cancelAppointmentButton := tgbotapi.NewInlineKeyboardButtonData(cancelText, fmt.Sprintf("cancel_%d", appointmentEntity.GetAppointmentId()))
		rescheduleButton := tgbotapi.NewInlineKeyboardButtonData(rescheduleText, fmt.Sprintf("reschedule_%d", appointmentEntity.GetAppointmentId()))
		keyboardRow := tgbotapi.NewInlineKeyboardRow(cancelAppointmentButton, rescheduleButton)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardRow)
		msg.ReplyMarkup = keyboard

	} else {
		emptyMessageText, err := c.messageService.GetMessage(ctx, userEntity, "empty appointments")
		if err != nil {
			msg = tgbotapi.NewMessage(c.tgUser.TgID, message.ServerError)
			c.bot.SendMessage(msg, messageDTO)
			return
		}
		msg = tgbotapi.NewMessage(c.tgUser.TgID, emptyMessageText)
	}

	c.machine.SetState(userEntity, userEntity.GetState(), state_machine.DetailMyAppointment)

	c.bot.SendMessage(msg, messageDTO)
}

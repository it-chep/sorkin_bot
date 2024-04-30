package start

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
	//_, _ = c.bot.Bot.Send(msg)
	appointmentId, err := strconv.Atoi(callbackData)
	if err != nil {
		msg = tgbotapi.NewMessage(c.tgUser.TgID, message.ServerError)
		_, _ = c.bot.Bot.Send(msg)
		return
	}

	appointmentEntity := c.appointmentService.GetAppointmentDetail(ctx, userEntity, appointmentId)

	if appointmentEntity.GetAppointmentId() != 0 {
		msg = tgbotapi.NewMessage(c.tgUser.TgID, fmt.Sprintf("Запись на %s %d", appointmentEntity.GetTimeStart(), appointmentEntity.GetAppointmentId()))
		cancelText, err := c.messageService.GetMessage(ctx, userEntity, "cancel appointment button")
		if err != nil {
			msg = tgbotapi.NewMessage(c.tgUser.TgID, message.ServerError)
			_, _ = c.bot.Bot.Send(msg)
			return
		}
		rescheduleText, err := c.messageService.GetMessage(ctx, userEntity, "reschedule appointment button")
		if err != nil {
			msg = tgbotapi.NewMessage(c.tgUser.TgID, message.ServerError)
			_, _ = c.bot.Bot.Send(msg)
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
			_, _ = c.bot.Bot.Send(msg)
			return
		}
		msg = tgbotapi.NewMessage(c.tgUser.TgID, emptyMessageText)
	}

	c.machine.SetState(userEntity, userEntity.GetState(), state_machine.ChooseSpeciality)

	sentMessage, err := c.bot.Bot.Send(msg)
	// todo мб вынести в отдельный метод
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s", err))
	}
	messageDTO.MessageID = int64(sentMessage.MessageID)
	messageDTO.Text = sentMessage.Text

	go func() {
		err := c.messageService.SaveMessageLog(context.TODO(), messageDTO)
		if err != nil {
			return
		}
	}()
}

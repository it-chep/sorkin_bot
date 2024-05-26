package callback

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/entity/user/state_machine"
)

func (c *CallbackBotMessage) mainMenu(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	switch callbackData {
	case "fast_appointment":
		c.botGateway.SendFastAppointmentMessage(ctx, userEntity, messageDTO)
		c.machine.SetState(userEntity, state_machine.FastAppointment)

	case "create_appointment":
		c.moreLessSpeciality(ctx, messageDTO, userEntity, callbackData)
		c.machine.SetState(userEntity, state_machine.ChooseSpeciality)
		c.appointmentService.CreateDraftAppointment(ctx, userEntity.GetTgId())

	case "my_appointments":
		messageId := c.botGateway.SendWaitMessage(ctx, userEntity, messageDTO, "wait my appointments")
		appointments := c.appointmentService.GetAppointments(ctx, userEntity)
		c.bot.RemoveMessage(userEntity.GetTgId(), messageId)

		if len(appointments) != 0 {
			c.botGateway.SendMyAppointmentsMessage(ctx, userEntity, appointments, messageDTO)
			c.machine.SetState(userEntity, state_machine.ChooseAppointment)
		} else {
			msgText, _ := c.messageService.GetMessage(ctx, userEntity, "empty appointments")
			msg := tgbotapi.NewMessage(userEntity.GetTgId(), msgText)
			c.bot.SendMessage(msg, messageDTO)

			c.botGateway.SendStartMessage(ctx, userEntity, messageDTO)
			c.machine.SetState(userEntity, state_machine.Start)
		}

	case "change_language":
		c.botGateway.SendChangeLanguageMessage(ctx, userEntity, messageDTO)
		c.machine.SetState(userEntity, state_machine.ChooseLanguage)
	}
}

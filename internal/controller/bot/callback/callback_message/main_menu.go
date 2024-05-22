package callback

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/entity/user/state_machine"
)

func (c *CallbackBotMessage) mainMenu(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	switch callbackData {
	case "fast_appointment":
		c.botGateway.SendFastAppointmentMessage(ctx, userEntity, messageDTO)
		go c.machine.SetState(userEntity, state_machine.FastAppointment)

	case "create_appointment":
		c.botGateway.SendChangeLanguageMessage(ctx, userEntity, messageDTO)
		go c.machine.SetState(userEntity, state_machine.ChooseSpeciality)

	case "my_appointments":
		appointments := c.appointmentService.GetAppointments(ctx, userEntity)
		if len(appointments) == 0 {
			c.botGateway.SendMyAppointmentsMessage(ctx, userEntity, appointments, messageDTO)
			go c.machine.SetState(userEntity, state_machine.MyAppointments)
		} else {
			c.botGateway.SendStartMessage(ctx, userEntity, messageDTO)
			go c.machine.SetState(userEntity, state_machine.Start)
		}

	case "change_language":
		c.botGateway.SendChangeLanguageMessage(ctx, userEntity, messageDTO)
		go c.machine.SetState(userEntity, state_machine.ChooseLanguage)
	}
}

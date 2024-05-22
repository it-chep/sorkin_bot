package bot_gateway

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
)

func (bg BotGateway) SendFastAppointmentMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	schedulesMap := bg.appointmentService.GetFastAppointmentSchedules(ctx)

	msgText, keyboard := bg.keyboard.ConfigureFastAppointmentMessage(ctx, user, schedulesMap)
	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	sentMessage := bg.bot.SendMessage(msg, messageDTO)
	go func() {
		err := bg.messageService.SaveMessageLog(ctx, sentMessage)
		if err != nil {
			bg.logger.Error(fmt.Sprintf("%s", err))
		}
	}()
}

func (bg BotGateway) SendConfirmAppointmentMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, doctorId int) {
	msgText, keyboard := bg.keyboard.ConfigureConfirmAppointmentMessage(ctx, user, doctorId)
	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	sentMessage := bg.bot.SendMessage(msg, messageDTO)
	go func() {
		err := bg.messageService.SaveMessageLog(ctx, sentMessage)
		if err != nil {
			bg.logger.Error(fmt.Sprintf("%s", err))
		}
	}()
}

func (bg BotGateway) SendMyAppointmentsMessage(ctx context.Context, user entity.User, appointments []appointment.Appointment, messageDTO tg.MessageDTO) {
	msgText, keyboard := bg.keyboard.ConfigureGetMyAppointmentsMessage(ctx, user, appointments, 0)
	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	sentMessage := bg.bot.SendMessage(msg, messageDTO)
	go func() {
		err := bg.messageService.SaveMessageLog(ctx, sentMessage)
		if err != nil {
			bg.logger.Error(fmt.Sprintf("%s", err))
		}
	}()
}

func (bg BotGateway) SendDetailAppointmentMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, appointmentId int) {

}

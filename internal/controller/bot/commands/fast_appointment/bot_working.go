package fast_appointment

import (
	"context"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/appointment"
	"sorkin_bot/internal/domain/services/message"
	"sorkin_bot/internal/domain/services/user"
	"sorkin_bot/pkg/client/telegram"
)

type FastAppointmentBotCommand struct {
	logger             *slog.Logger
	bot                telegram.Bot
	tgUser             tg.TgUserDTO
	userService        user.UserService
	machine            *state_machine.UserStateMachine
	appointmentService appointment.AppointmentService
	messageService     message.MessageService
}

func NewFastAppointmentBotCommand(
	logger *slog.Logger,
	bot telegram.Bot,
	tgUser tg.TgUserDTO,
	userService user.UserService,
	machine *state_machine.UserStateMachine,
	appointmentService appointment.AppointmentService,
	messageService message.MessageService,
) FastAppointmentBotCommand {
	return FastAppointmentBotCommand{
		logger:             logger,
		bot:                bot,
		tgUser:             tgUser,
		userService:        userService,
		machine:            machine,
		appointmentService: appointmentService,
		messageService:     messageService,
	}
}

func (c *FastAppointmentBotCommand) Execute(ctx context.Context, message tg.MessageDTO) {
	//c.appointmentService.Mis.FastAppointment(ctx)
	//var _ tgbotapi.MessageConfig
	//if err != nil {
	//	return
	//}
	//c.machine.SetState(user, user.GetState(), state_machine.FastAppointment)
	//
	//msg = tgbotapi.NewMessage(c.tgUser.TgID, "FastAppointmentBotCommand message")
	//c.logger.Info(fmt.Sprintf("%s", message))
	//
	//sentMessage, err := c.bot.Bot.Send(msg)
	//// todo мб вынести в отдельный метод
	//if err != nil {
	//	c.logger.Error(fmt.Sprintf("%s", err))
	//}
	//message.MessageID = int64(sentMessage.MessageID)
	//message.Text = sentMessage.Text
	//
	//go func() {
	//	err := c.messageService.SaveMessageLog(context.TODO(), message)
	//	if err != nil {
	//		return
	//	}
	//}()
}

package cancel_appointment

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/appointment"
	"sorkin_bot/internal/domain/services/message"
	"sorkin_bot/internal/domain/services/user"
	"sorkin_bot/pkg/client/telegram"
)

type CancelAppointmentBotCommand struct {
	logger             *slog.Logger
	bot                telegram.Bot
	tgUser             tg.TgUserDTO
	userService        user.UserService
	machine            *state_machine.UserStateMachine
	appointmentService appointment.AppointmentService
	messageService     message.MessageService
}

func NewCancelAppointmentBotCommand(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO, userService user.UserService, machine *state_machine.UserStateMachine, appointmentService appointment.AppointmentService, messageService message.MessageService,
) CancelAppointmentBotCommand {
	return CancelAppointmentBotCommand{
		logger:             logger,
		bot:                bot,
		tgUser:             tgUser,
		userService:        userService,
		machine:            machine,
		appointmentService: appointmentService,
		messageService:     messageService,
	}
}

// Execute место связи telegram и бизнес логи
func (c *CancelAppointmentBotCommand) Execute(ctx context.Context, message tg.MessageDTO) {
	var msg tgbotapi.MessageConfig
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser)
	if userEntity.GetState() != "" {
		err, appointments := c.appointmentService.MyAppointments(ctx)
		if err != nil {
			return
		}
		msg = tgbotapi.NewMessage(c.tgUser.TgID, "Please select appointment")
		keyboard := tgbotapi.NewInlineKeyboardMarkup()

		for _, appointmentEntity := range appointments.Data {
			btn := tgbotapi.NewInlineKeyboardButtonData(appointmentEntity.TimeStart, fmt.Sprintf("%d", appointmentEntity.Id))
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
		msg.ReplyMarkup = keyboard

	} else {
		return
	}

	c.machine.SetState(userEntity, userEntity.GetState(), state_machine.ChooseAppointment)

	//msg = tgbotapi.NewMessage(c.tgUser.TgID, "Выберите запись, которую хотите отменить")
	c.logger.Info(fmt.Sprintf("%s", message))

	sentMessage, err := c.bot.Bot.Send(msg)
	// todo мб вынести в отдельный метод
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s", err))
	}
	message.MessageID = int64(sentMessage.MessageID)
	message.Text = sentMessage.Text

	go func() {
		err := c.messageService.SaveMessageLog(context.TODO(), message)
		if err != nil {
			return
		}
	}()
}

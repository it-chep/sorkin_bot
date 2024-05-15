package cancel_appointment

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/pkg/client/telegram"
)

type CancelAppointmentBotCommand struct {
	logger             *slog.Logger
	bot                telegram.Bot
	tgUser             tg.TgUserDTO
	userService        UserService
	machine            *state_machine.UserStateMachine
	appointmentService AppointmentService
	messageService     MessageService
}

func NewCancelAppointmentBotCommand(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO, userService UserService, machine *state_machine.UserStateMachine, appointmentService AppointmentService, messageService MessageService,
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
	// так как мы не изменяем бизнес сущность, а бот меняет состояние, то нахождение сущность в слое controllers некритично
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser)

	if userEntity.GetState() != "" {
		appointments := c.appointmentService.GetAppointments(ctx, userEntity)
		messageText, err := c.messageService.GetMessage(ctx, userEntity, "Select appointment")
		msg = tgbotapi.NewMessage(c.tgUser.TgID, messageText)
		if err != nil {
			c.bot.SendMessage(msg, message)
			return
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup()

		for _, appointmentEntity := range appointments {
			btn := tgbotapi.NewInlineKeyboardButtonData(appointmentEntity.GetTimeStart(), fmt.Sprintf("%d", appointmentEntity.GetAppointmentId()))
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
		msg.ReplyMarkup = keyboard

	} else {
		return
	}

	c.machine.SetState(userEntity, userEntity.GetState(), state_machine.ChooseAppointment)

	c.bot.SendMessage(msg, messageDTO)
}

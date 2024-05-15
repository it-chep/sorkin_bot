package my_appointment

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/pkg/client/telegram"
)

type MyAppointmentsCommand struct {
	logger             *slog.Logger
	bot                telegram.Bot
	tgUser             tg.TgUserDTO
	machine            *state_machine.UserStateMachine
	userService        UserService
	appointmentService AppointmentService
	messageService     MessageService
}

func NewMyAppointmentsCommand(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO, machine *state_machine.UserStateMachine, userService UserService, appointmentService AppointmentService, messageService MessageService) MyAppointmentsCommand {
	return MyAppointmentsCommand{
		logger:             logger,
		bot:                bot,
		tgUser:             tgUser,
		machine:            machine,
		userService:        userService,
		appointmentService: appointmentService,
		messageService:     messageService,
	}
}

func (c MyAppointmentsCommand) Execute(ctx context.Context, messageDTO tg.MessageDTO) {
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser)
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	msg := tgbotapi.NewMessage(c.tgUser.TgID, "Выберите запись")

	appointments := c.appointmentService.GetAppointments(ctx, userEntity)
	for _, appointmentEntity := range appointments {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s - %s", appointmentEntity.GetTimeStart(), appointmentEntity.GetTimeEnd()),
			fmt.Sprintf("%d", appointmentEntity.GetAppointmentId()),
		)
		row := tgbotapi.NewInlineKeyboardRow(btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	msg.ReplyMarkup = keyboard

	c.machine.SetState(userEntity, userEntity.GetState(), state_machine.ChooseAppointment)

	c.bot.SendMessage(msg, messageDTO)
}

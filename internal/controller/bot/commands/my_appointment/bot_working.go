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
	userService        userService
	appointmentService appointmentService
	messageService     messageService
}

func NewMyAppointmentsCommand(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO, machine *state_machine.UserStateMachine, userService userService, appointmentService appointmentService, messageService messageService) MyAppointmentsCommand {
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
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser.TgID)
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	msg := tgbotapi.NewMessage(c.tgUser.TgID, "Выберите запись")

	appointments := c.appointmentService.GetAppointments(ctx, userEntity)

	if len(appointments) == 0 {
		c.machine.SetState(userEntity, *userEntity.GetState(), state_machine.Start)
		msgText, _ := c.messageService.GetMessage(ctx, userEntity, "Start")
		msg := tgbotapi.NewMessage(c.tgUser.TgID, msgText)
		c.bot.SendMessage(msg, messageDTO)
		return
	}

	for _, appointmentEntity := range appointments {
		appointmentEntity.GetTimeStart()
		buttonDoc, _ := c.messageService.GetMessage(ctx, userEntity, "doc information button")
		docBtn := tgbotapi.NewInlineKeyboardButtonData(buttonDoc, fmt.Sprintf("doc_info_%d", appointmentEntity.GetDoctorId()))

		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s: %s - %s", appointmentEntity.GetDate(), appointmentEntity.GetTimeStartShort(), appointmentEntity.GetTimeEndShort()),
			fmt.Sprintf("appointmentId_%d", appointmentEntity.GetAppointmentId()),
		)

		row := tgbotapi.NewInlineKeyboardRow(docBtn, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	msg.ReplyMarkup = keyboard
	c.machine.SetState(userEntity, *userEntity.GetState(), state_machine.ChooseAppointment)
	c.bot.SendMessage(msg, messageDTO)
}

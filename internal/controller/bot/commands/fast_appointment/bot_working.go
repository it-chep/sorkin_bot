package fast_appointment

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/pkg/client/telegram"
)

type FastAppointmentBotCommand struct {
	logger             *slog.Logger
	bot                telegram.Bot
	tgUser             tg.TgUserDTO
	userService        userService
	machine            *state_machine.UserStateMachine
	appointmentService appointmentService
	messageService     messageService
	botService         botService
}

func NewFastAppointmentBotCommand(
	logger *slog.Logger,
	bot telegram.Bot,
	tgUser tg.TgUserDTO,
	userService userService,
	machine *state_machine.UserStateMachine,
	appointmentService appointmentService,
	messageService messageService,
	botService botService,
) FastAppointmentBotCommand {
	return FastAppointmentBotCommand{
		logger:             logger,
		bot:                bot,
		tgUser:             tgUser,
		userService:        userService,
		machine:            machine,
		appointmentService: appointmentService,
		messageService:     messageService,
		botService:         botService,
	}
}

func (c *FastAppointmentBotCommand) Execute(ctx context.Context, message tg.MessageDTO) {
	var msg tgbotapi.MessageConfig
	userEntity, err := c.userService.GetUser(ctx, c.tgUser.TgID)

	if err != nil {
		return
	}

	schedulesMap := c.appointmentService.GetFastAppointmentSchedules(ctx)

	msgText, keyboard := c.botService.ConfigureFastAppointmentMessage(ctx, userEntity, schedulesMap)
	msg = tgbotapi.NewMessage(c.tgUser.TgID, msgText)
	msg.ReplyMarkup = keyboard
	c.bot.SendMessage(msg, message)

	go c.machine.SetState(userEntity, *userEntity.GetState(), state_machine.FastAppointment)
	go c.appointmentService.CreateDraftAppointment(ctx, userEntity.GetTgId())
}

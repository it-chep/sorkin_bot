package fast_appointment

import (
	"context"
	"fmt"
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
}

func NewFastAppointmentBotCommand(
	logger *slog.Logger,
	bot telegram.Bot,
	tgUser tg.TgUserDTO,
	userService userService,
	machine *state_machine.UserStateMachine,
	appointmentService appointmentService,
	messageService messageService,
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
	var msg tgbotapi.MessageConfig

	schedulesMap := c.appointmentService.GetFastAppointmentSchedules(ctx)

	var rows [][]tgbotapi.InlineKeyboardButton

	for doctorId, schedules := range schedulesMap {
		if len(schedules) > 1 {
			for _, schedule := range schedules {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("doctorId_%d__timeStart_%s__timeEnd_%s", doctorId, schedule.GetTimeStart(), schedule.GetTimeEnd()),
					fmt.Sprintf("doctorId_%d__timeStart_'%s'__timeEnd_'%s'", doctorId, schedule.GetTimeStart(), schedule.GetTimeEnd()),
				)))
			}
		} else if len(schedules) == 1 {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("doctorId_%d__timeStart_%s__timeEnd_%s", doctorId, schedules[0].GetTimeStart(), schedules[0].GetTimeEnd()),
				fmt.Sprintf("doctorId_%d__timeStart_'%s'__timeEnd_'%s'", doctorId, schedules[0].GetTimeStart(), schedules[0].GetTimeEnd()),
			)))
		} else {
			return
		}
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg = tgbotapi.NewMessage(c.tgUser.TgID, "FastAppointmentBotCommand message")
	msg.ReplyMarkup = keyboard

	c.logger.Info(fmt.Sprintf("%s", message))

	c.bot.SendMessage(msg, message)

	// todo, мб горутину на стейты
	userEntity, err := c.userService.GetUser(ctx, c.tgUser.TgID)
	if err != nil {
		return
	}
	go c.machine.SetState(userEntity, userEntity.GetState(), state_machine.FastAppointment)
	go c.appointmentService.CreateDraftAppointment(ctx, userEntity.GetTgId())
}

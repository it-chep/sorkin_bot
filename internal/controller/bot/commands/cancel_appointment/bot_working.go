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
	// так как мы не изменяем бизнес сущность, а бот меняет состояние, то нахождение сущность в слое controllers некритично
	userEntity, _ := c.userService.GetUser(ctx, c.tgUser)

	if userEntity.GetState() != "" {
		err, appointments := c.appointmentService.Mis.MyAppointments(ctx)
		if err != nil {
			return
		}
		messageText, err := c.messageService.GetMessage(ctx, userEntity, "Select appointment")
		msg = tgbotapi.NewMessage(c.tgUser.TgID, messageText)
		if err != nil {
			_, _ = c.bot.Bot.Send(msg)
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
	//todo докрутить логику со специальностями
	if userEntity.GetState() == state_machine.ChooseSpeciality {
		err, specialities := c.appointmentService.Mis.GetSpecialities(ctx)
		if err != nil {
			return
		}
		messageText, err := c.messageService.GetMessage(ctx, userEntity, "Choose speciality")
		msg = tgbotapi.NewMessage(c.tgUser.TgID, messageText)
		if err != nil {
			_, _ = c.bot.Bot.Send(msg)
			return
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup()
		translatedSpecialities, err := c.appointmentService.GetTranslatedSpecialities(ctx, userEntity, specialities)
		if err != nil {
			return
		}
		for speciallityId, translatedSpeciality := range translatedSpecialities {
			btn := tgbotapi.NewInlineKeyboardButtonData(translatedSpeciality, fmt.Sprintf("%d", speciallityId))
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
		msg.ReplyMarkup = keyboard

	} else {
		return
	}
	//todo докрутить логику со специальностями

	c.machine.SetState(userEntity, userEntity.GetState(), state_machine.ChooseAppointment)

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

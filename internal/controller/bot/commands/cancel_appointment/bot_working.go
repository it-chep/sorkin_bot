package cancel_appointment

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/message"
	"sorkin_bot/internal/domain/services/user"
	"sorkin_bot/pkg/client/telegram"
)

type CancelAppointmentBotCommand struct {
	logger         *slog.Logger
	bot            telegram.Bot
	tgUser         tg.TgUserDTO
	userService    user.UserService
	machine        *state_machine.UserStateMachine
	messageService message.MessageService
}

func NewCancelAppointmentBotCommand(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO, userService user.UserService, machine *state_machine.UserStateMachine, messageService message.MessageService,
) CancelAppointmentBotCommand {
	return CancelAppointmentBotCommand{
		logger:         logger,
		bot:            bot,
		tgUser:         tgUser,
		userService:    userService,
		machine:        machine,
		messageService: messageService,
	}
}

// Execute место связи telegram и бизнес логи
func (c *CancelAppointmentBotCommand) Execute(ctx context.Context, message tg.MessageDTO) {
	user, err := c.userService.CancelAppointment(ctx, c.tgUser)
	var msg tgbotapi.MessageConfig

	if err != nil {
		return
	}
	c.machine.SetState(user, user.GetState(), "cancelAppointment")

	msg = tgbotapi.NewMessage(c.tgUser.TgID, "CancelAppointmentBotCommand message")
	c.logger.Info(fmt.Sprintf("%s", message))

	sentMessage, err := c.bot.Bot.Send(msg)
	message.MessageID = int64(sentMessage.MessageID)
	message.Text = sentMessage.Text

	if err != nil {
		c.logger.Error(fmt.Sprintf("%s", err))
	}
	go func() {
		err := c.messageService.SaveMessageLog(context.TODO(), message)
		if err != nil {
			return
		}
	}()
}

package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log/slog"
	"net/http"
	"sorkin_bot/internal/config"
	callback "sorkin_bot/internal/controller/bot/callback/callback_message"
	"sorkin_bot/internal/controller/bot/text/text_message"

	"sorkin_bot/internal/controller/bot/commands/cancel_appointment"
	"sorkin_bot/internal/controller/bot/commands/change_language"
	"sorkin_bot/internal/controller/bot/commands/create_appointment"
	"sorkin_bot/internal/controller/bot/commands/fast_appointment"
	"sorkin_bot/internal/controller/bot/commands/my_appointment"
	"sorkin_bot/internal/controller/bot/commands/start"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/pkg/client/telegram"
)

type TelegramWebhookController struct {
	cfg                config.Config
	logger             *slog.Logger
	bot                telegram.Bot
	machine            *state_machine.UserStateMachine
	userService        UserService
	appointmentService AppointmentService
	messageService     MessageService
	botService         BotService
}

func NewTelegramWebhookController(
	cfg config.Config,
	logger *slog.Logger,
	bot telegram.Bot,
	machine *state_machine.UserStateMachine,
	userService UserService,
	appointmentService AppointmentService,
	messageService MessageService,
	botService BotService,
) TelegramWebhookController {

	return TelegramWebhookController{
		cfg:                cfg,
		logger:             logger,
		bot:                bot,
		machine:            machine,
		userService:        userService,
		appointmentService: appointmentService,
		messageService:     messageService,
		botService:         botService,
	}
}

func (t TelegramWebhookController) BotWebhookHandler(c *gin.Context) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.logger.Error(fmt.Sprintf("%s", err))
		}
	}(c.Request.Body)

	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		t.logger.Error("Error binding JSON", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	tgUser := t.getUserFromWebhook(update)
	tgMessage := t.getMessageFromWebhook(update)
	// Сначала проверяем на команду, потом на текстовое сообщение, потом callback
	if update.Message != nil {
		ctx := context.WithValue(context.Background(), "userID", update.Message.From.ID)
		if update.Message.IsCommand() {
			t.ForkCommands(ctx, update, tgUser, tgMessage)
		} else {
			t.ForkMessages(ctx, tgUser, tgMessage)
		}
	} else if update.CallbackQuery != nil {
		ctx := context.WithValue(context.Background(), "userID", update.CallbackQuery.From.ID)
		t.ForkCallbacks(ctx, update, tgUser, tgMessage)
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (t TelegramWebhookController) ForkCommands(ctx context.Context, update tgbotapi.Update, tgUser tg.TgUserDTO, tgMessage tg.MessageDTO) {

	switch update.Message.Command() {
	case "start":
		command := start.NewStartBotCommand(t.logger, t.bot, tgUser, t.userService, t.messageService, t.botService)
		command.Execute(ctx, tgMessage)
	case "help":
		// service по работе с help
		t.bot.SendMessage(tgbotapi.NewMessage(update.FromChat().ID, "i will help you"), tgMessage)
	case "tech_support":
		// service по работе с tech_support
		t.bot.SendMessage(tgbotapi.NewMessage(update.FromChat().ID, "i will help you"), tgMessage)
	case "fast_appointment":
		// service по работе с fast appointment
		command := fast_appointment.NewFastAppointmentBotCommand(t.logger, t.bot, tgUser, t.userService, t.machine, t.appointmentService, t.messageService)
		command.Execute(ctx, tgMessage)
	case "appointment":
		// service по работе с appointment
		command := create_appointment.NewCreateAppointmentCommand(t.logger, t.bot, tgUser, t.userService, t.machine, t.appointmentService, t.messageService, t.botService)
		command.Execute(ctx, tgMessage)
	case "cancel_appointment":
		// service по работе с cancel_appointment
		command := cancel_appointment.NewCancelAppointmentBotCommand(t.logger, t.bot, tgUser, t.userService, t.machine, t.appointmentService, t.messageService)
		command.Execute(ctx, tgMessage)
	case "my_appointments":
		// service по работе с my_appointments
		command := my_appointment.NewMyAppointmentsCommand(t.logger, t.bot, tgUser, t.machine, t.userService, t.appointmentService, t.messageService)
		command.Execute(ctx, tgMessage)
	case "change_language":
		command := change_language.NewChangeLanguageCommand(t.logger, t.bot, tgUser, t.userService, t.messageService, t.botService)
		command.Execute(ctx, tgMessage)
	case "menu":
		// service по работе с menu
		t.bot.SendMessage(tgbotapi.NewMessage(update.FromChat().ID, "i will help you"), tgMessage)
	}
}

// todo в эти форки будут сыпаться все текстовые сообщения и колбэки
// todo и в зависимости от состояния пользователя ему будет выдаваться контент

func (t TelegramWebhookController) ForkMessages(ctx context.Context, tgUser tg.TgUserDTO, tgMessage tg.MessageDTO) {
	messageBot := text_message.NewTextBotMessage(t.logger, t.bot, tgUser, t.machine, t.userService, t.messageService, t.appointmentService, t.botService)
	messageBot.Execute(ctx, tgMessage)
}

func (t TelegramWebhookController) ForkCallbacks(ctx context.Context, update tgbotapi.Update, tgUser tg.TgUserDTO, tgMessage tg.MessageDTO) {
	callbackData := update.CallbackData()
	callbackBot := callback.NewCallbackBot(t.logger, t.bot, tgUser, t.machine, t.userService, t.messageService, t.appointmentService, t.botService)
	callbackBot.Execute(ctx, tgMessage, callbackData)
}

func (t TelegramWebhookController) getUserFromWebhook(update tgbotapi.Update) tg.TgUserDTO {
	var tgUser tg.TgUserDTO
	var userJSON []byte
	var err error

	// Todo возможно улучшить
	if update.CallbackQuery != nil {
		userJSON, err = json.Marshal(update.CallbackQuery.From)
	} else {
		userJSON, err = json.Marshal(update.Message.From)
	}

	if err != nil {
		t.logger.Error(fmt.Sprintf("Error marshaling user to JSON: %s", err))
		return tg.TgUserDTO{}
	}
	if err = json.Unmarshal(userJSON, &tgUser); err != nil {
		t.logger.Error(fmt.Sprintf("Error decoding JSON: %s", err))
		return tg.TgUserDTO{}
	}
	return tgUser
}

func (t TelegramWebhookController) getMessageFromWebhook(update tgbotapi.Update) tg.MessageDTO {
	var tgMessage tg.MessageDTO
	var userJSON []byte
	var err error

	// Todo возможно улучшить
	if update.CallbackQuery != nil {
		userJSON, err = json.Marshal(update.CallbackQuery.Message)
	} else {
		userJSON, err = json.Marshal(update.Message)
	}

	if err != nil {
		t.logger.Error(fmt.Sprintf("Error marshaling user to JSON: %s", err))
		return tg.MessageDTO{}
	}
	if err = json.Unmarshal(userJSON, &tgMessage); err != nil {
		t.logger.Error(fmt.Sprintf("Error decoding JSON: %s", err))
		return tg.MessageDTO{}
	}
	return tgMessage
}

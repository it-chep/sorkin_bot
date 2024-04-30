package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log/slog"
	"net/http"
	"sorkin_bot/internal/config"
	start2 "sorkin_bot/internal/controller/bot/callback/callback_message"
	"sorkin_bot/internal/controller/bot/commands/cancel_appointment"
	"sorkin_bot/internal/controller/bot/commands/change_language"
	"sorkin_bot/internal/controller/bot/commands/create_appointment"
	"sorkin_bot/internal/controller/bot/commands/fast_appointment"
	"sorkin_bot/internal/controller/bot/commands/my_appointment"
	"sorkin_bot/internal/controller/bot/commands/start"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"sorkin_bot/internal/domain/services/appointment"
	"sorkin_bot/internal/domain/services/bot"
	"sorkin_bot/internal/domain/services/message"
	"sorkin_bot/internal/domain/services/user"
	"sorkin_bot/pkg/client/telegram"
)

type TelegramWebhookController struct {
	router             *gin.Engine
	cfg                config.Config
	logger             *slog.Logger
	bot                telegram.Bot
	machine            *state_machine.UserStateMachine
	userService        user.UserService
	appointmentService appointment.AppointmentService
	messageService     message.MessageService
	botService         bot.BotService
}

func NewTelegramWebhookController(
	cfg config.Config,
	logger *slog.Logger,
	bot telegram.Bot,
	machine *state_machine.UserStateMachine,
	userService user.UserService,
	appointmentService appointment.AppointmentService,
	messageService message.MessageService,
	botService bot.BotService,
) TelegramWebhookController {
	router := gin.New()
	router.Use(gin.Recovery())

	return TelegramWebhookController{
		router:             router,
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

	// Сначала проверяем на команду, потом на текстовое сообщение, потом callback
	if update.Message != nil {
		if update.Message.IsCommand() {
			err := t.ForkCommands(update)
			if err != nil {
				return
			}
		} else {
			err := t.ForkMessages(update)
			if err != nil {
				return
			}
		}
	} else if update.CallbackQuery != nil {
		err := t.ForkCallbacks(update)
		if err != nil {
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (t TelegramWebhookController) ForkCommands(update tgbotapi.Update) error {
	tgUser := t.getUserFromWebhook(update)
	tgMessage := t.getMessageFromWebhook(update)
	ctx := context.WithValue(context.Background(), "userID", update.Message.From.ID)

	switch update.Message.Command() {
	case "start":
		t.logger.Info("start command was called")
		command := start.NewStartBotCommand(t.logger, t.bot, tgUser, t.userService, t.messageService, t.botService)
		command.Execute(ctx, tgMessage)
	case "help":
		// service по работе с help

		_, err := t.bot.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "i will help you"))
		if err != nil {
			return err
		}
		return nil
	case "tech_support":
		// service по работе с tech_support

		_, err := t.bot.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "i will help you"))
		if err != nil {
			return err
		}
		return nil
	case "fast_appointment":
		// service по работе с fast appointment

		t.logger.Info("fast_appointment command was called")
		command := fast_appointment.NewFastAppointmentBotCommand(t.logger, t.bot, tgUser, t.userService, t.machine, t.appointmentService, t.messageService)
		command.Execute(ctx, tgMessage)
		return nil
	case "appointment":
		// service по работе с appointment
		t.logger.Info("create_appointment command was called")

		command := create_appointment.NewCreateAppointmentCommand(t.logger, t.bot, tgUser, t.userService, t.machine, t.appointmentService, t.messageService)
		command.Execute(ctx, tgMessage)
	case "cancel_appointment":
		// service по работе с cancel_appointment

		t.logger.Info("cancel_appointment command was called")
		command := cancel_appointment.NewCancelAppointmentBotCommand(t.logger, t.bot, tgUser, t.userService, t.machine, t.appointmentService, t.messageService)
		command.Execute(ctx, tgMessage)
		return nil
	case "my_appointments":
		// service по работе с my_appointments

		command := my_appointment.NewMyAppointmentsCommand(t.logger, t.bot, tgUser, t.machine, t.userService, t.appointmentService, t.messageService)
		command.Execute(ctx, tgMessage)
		return nil
	case "change_language":
		t.logger.Info("change_language command was called")

		command := change_language.NewChangeLanguageCommand(t.logger, t.bot, tgUser, t.userService, t.messageService, t.botService)
		command.Execute(ctx)
	case "menu":
		// service по работе с menu

		_, err := t.bot.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "my_appointments"))
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("no commands")
}

// todo в эти форки будут сыпаться все текстовые сообщения и колбэки
// todo и в зависимости от состояния пользователя ему будет выдаваться контент

func (t TelegramWebhookController) ForkMessages(update tgbotapi.Update) error {
	tgMessage := t.getMessageFromWebhook(update)
	t.logger.Info(tgMessage.Text)
	return errors.New("no texts yet")
}

func (t TelegramWebhookController) ForkCallbacks(update tgbotapi.Update) error {
	ctx := context.TODO()
	tgUser := t.getUserFromWebhook(update)
	callbackData := update.CallbackData()
	tgMessage := t.getMessageFromWebhook(update)
	callback := start2.NewCallbackBot(t.logger, t.bot, tgUser, t.machine, t.userService, t.messageService, t.appointmentService)
	callback.Execute(ctx, tgMessage, callbackData)
	return errors.New("no callbacks yet")
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
	if err := json.Unmarshal(userJSON, &tgUser); err != nil {
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
	if err := json.Unmarshal(userJSON, &tgMessage); err != nil {
		t.logger.Error(fmt.Sprintf("Error decoding JSON: %s", err))
		return tg.MessageDTO{}
	}

	go func() {
		err := t.messageService.SaveMessageLog(context.TODO(), tgMessage)
		if err != nil {
			return
		}
	}()
	if err != nil {
		t.logger.Error(fmt.Sprintf("Error while saving message: %s", err))
		return tg.MessageDTO{}
	}

	return tgMessage
}

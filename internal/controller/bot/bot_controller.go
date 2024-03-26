package bot

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log/slog"
	"net/http"
	"sorkin_bot/internal/config"
	"sorkin_bot/pkg/client/telegram"
)

type TelegramWebhookController struct {
	router *gin.Engine
	cfg    config.Config
	logger *slog.Logger
	bot    telegram.Bot
}

func NewTelegramWebhookController(cfg config.Config, logger *slog.Logger, bot telegram.Bot) TelegramWebhookController {
	router := gin.New()
	router.Use(gin.Recovery())

	return TelegramWebhookController{
		router: router,
		cfg:    cfg,
		logger: logger,
		bot:    bot,
	}
}

func (t TelegramWebhookController) BotWebhookHandler(c *gin.Context) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.logger.Error(fmt.Sprintf(""))
		}
	}(c.Request.Body)

	t.logger.Info("Received webhook update")

	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		t.logger.Error("Error binding JSON", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if update.Message.IsCommand() {
		err := t.ForkCommands(update)
		if err != nil {
			return
		}
	}
	// Echobot
	//_, err := t.bot.Bot.Send(tgbotapi.NewMessage(update.Message.From.ID, update.Message.Text))
	//if err != nil {
	//	return
	//}
	// Echobot

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (t TelegramWebhookController) ForkCommands(update tgbotapi.Update) error {
	switch update.Message.Command() {
	case "start":
		_, err := t.bot.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "hello"))
		if err != nil {
			return err
		}
		return nil
	case "help":
		_, err := t.bot.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "i will help you"))
		if err != nil {
			return err
		}
		return nil
	case "fast_appointment":
		_, err := t.bot.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "really_fast"))
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("no commands")
}

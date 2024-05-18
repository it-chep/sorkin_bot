package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/config"
	"sorkin_bot/internal/controller/dto/tg"
)

type Bot struct {
	Bot            *tgbotapi.BotAPI
	logger         *slog.Logger
	messageService MessageService
}

func NewTelegramBot(cfg config.Config, logger *slog.Logger, messageService MessageService) *Bot {
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	bot.Debug = true
	if err != nil {
		panic("can't create bot instance")
	}
	//updates, err := bot.GetUpdates(tgbotapi.UpdateConfig{AllowedUpdates: []string{"message", "callback_query"}, Limit: 100, Offset: 0})
	//if err != nil {
	//	return nil
	//}
	//logger.Info(fmt.Sprintf("%v", updates))
	//time.Sleep(6 * time.Second)
	wh, _ := tgbotapi.NewWebhook(cfg.Bot.WebhookURL + bot.Token + "/")
	_, err = bot.Request(wh)
	if err != nil {
		panic("can't while request set webhook")
	}

	_, err = bot.GetWebhookInfo()

	if err != nil {
		panic("error while getting webhook")
	}
	return &Bot{
		Bot:            bot,
		logger:         logger,
		messageService: messageService,
	}
}

func (bot *Bot) SendMessage(msg tgbotapi.MessageConfig, messageDTO tg.MessageDTO) {
	sentMessage, err := bot.Bot.Send(msg)
	if err != nil {
		bot.logger.Error(fmt.Sprintf("%s", err))
	}

	messageDTO.MessageID = int64(sentMessage.MessageID)
	messageDTO.Text = sentMessage.Text

	go func() {
		err = bot.messageService.SaveMessageLog(context.Background(), messageDTO)
		if err != nil {
			bot.logger.Error(fmt.Sprintf("%s", err))
		}
	}()
}

func (bot *Bot) RemoveMessage(chatId int64, messageId int) {
	messageToDelete := tgbotapi.NewDeleteMessage(chatId, messageId)
	_, err := bot.Bot.Send(messageToDelete)
	if err != nil {
		return
	}
}

package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/services/message"
)

type BotService struct {
	logger         *slog.Logger
	messageService message.MessageService
}

func NewBotService(logger *slog.Logger, messageService message.MessageService) BotService {
	return BotService{
		//readRepo:                 readRepo,
		//administratorHelpUseCase: administratorHelpUseCase,
		logger:         logger,
		messageService: messageService,
	}
}

func (bs BotService) AdministratorHelp() {
	// 	TODO create request message_log - controller or adapter mb
	//  TODO get language
	//	TODO get admin message by language
	//  return message
	//	TODO create response message_log may be go send_message() {} - controller or adapter mb
}

func (bs BotService) ConfigureChangeLanguageMessage(ctx context.Context, user entity.User) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 EN", "EN"),
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 RU", "RU"),
			tgbotapi.NewInlineKeyboardButtonData("🇵🇹 PT", "PT"),
		),
	)
	msgText, _ = bs.messageService.GetMessage(ctx, user, "change_language")
	return msgText, keyboard
}

func (bs BotService) ConfigureGetPhoneMessage(ctx context.Context, user entity.User) (msgText string, keyboard tgbotapi.ReplyKeyboardMarkup) {
	buttonText, _ := bs.messageService.GetMessage(ctx, user, "send phone button")
	keyboard = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact(buttonText),
		),
	)
	msgText, _ = bs.messageService.GetMessage(ctx, user, "send phone message")
	return msgText, keyboard
}

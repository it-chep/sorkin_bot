package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/services/message"
)

type BotService struct {
	//readRepo                 ReadMessagesRepo
	//administratorHelpUseCase AdministratorHelpUseCase
	logger         *slog.Logger
	messageService message.MessageService
}

func NewBotService(logger *slog.Logger, messageService message.MessageService) BotService {
	return BotService{
		//readRepo:                 readRepo,
		//administratorHelpUseCase: administratorHelpUseCase,
		// MORE usecases
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

func (bs BotService) CancelAppointment() {
	// 	TODO create request message_log - controller or adapter mb

	//  TODO POST TO cancel_appointment
	//  TODO get language
	//	TODO get cancel_appointment message by language and by status
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

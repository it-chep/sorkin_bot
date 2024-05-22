package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	entity "sorkin_bot/internal/domain/entity/user"
)

func (bs BotService) ConfigureChangeLanguageMessage(ctx context.Context, userEntity entity.User) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 EN", "EN"),
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 RU", "RU"),
			tgbotapi.NewInlineKeyboardButtonData("🇵🇹 PT", "PT"),
		),
	)
	msgText, _ = bs.messageService.GetMessage(ctx, userEntity, "change_language")
	return msgText, keyboard
}

func (bs BotService) ConfigureGetPhoneMessage(ctx context.Context, userEntity entity.User) (msgText string, keyboard tgbotapi.ReplyKeyboardMarkup) {
	buttonText, _ := bs.messageService.GetMessage(ctx, userEntity, "send phone button")
	keyboard = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact(buttonText),
		),
	)
	msgText, _ = bs.messageService.GetMessage(ctx, userEntity, "enter phone")
	return msgText, keyboard
}

func (bs BotService) ConfigureMainMenuMessage(ctx context.Context, userEntity entity.User) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	buttonFastAppointment, _ := bs.messageService.GetMessage(ctx, userEntity, "btn start fast appointment")
	buttonCreateAppointment, _ := bs.messageService.GetMessage(ctx, userEntity, "btn start create appointment")
	buttonMyAppointments, _ := bs.messageService.GetMessage(ctx, userEntity, "btn my appointments")
	buttonChangeLanguage, _ := bs.messageService.GetMessage(ctx, userEntity, "btn choose language")

	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonFastAppointment, "fast_appointment"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonCreateAppointment, "create_appointment"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonMyAppointments, "my_appointments"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonChangeLanguage, "change_language"),
		),
	)

	msgText, _ = bs.messageService.GetMessage(ctx, userEntity, "Start")
	return msgText, keyboard
}

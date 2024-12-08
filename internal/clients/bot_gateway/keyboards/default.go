package keyboards

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	entity "sorkin_bot/internal/domain/entity/user"
)

func (k Keyboards) ConfigureChangeLanguageMessage(ctx context.Context, userEntity entity.User) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 EN", "EN"),
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 RU", "RU"),
			tgbotapi.NewInlineKeyboardButtonData("🇵🇹 PT", "PT"),
		),
	)
	msgText, _ = k.messageService.GetMessage(ctx, userEntity, "change_language")
	return msgText, keyboard
}

func (k Keyboards) ConfigureGetPhoneMessage(ctx context.Context, userEntity entity.User) (msgText string, keyboard tgbotapi.ReplyKeyboardMarkup) {
	buttonText, _ := k.messageService.GetMessage(ctx, userEntity, "send phone button")
	keyboard = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact(buttonText),
		),
	)
	msgText, _ = k.messageService.GetMessage(ctx, userEntity, "enter phone")
	return msgText, keyboard
}

func (k Keyboards) ConfigureMainMenuMessage(ctx context.Context, userEntity entity.User) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	buttonCreateAppointment, _ := k.messageService.GetMessage(ctx, userEntity, "btn start create appointment")
	buttonMyAppointments, _ := k.messageService.GetMessage(ctx, userEntity, "btn my appointments")
	buttonChangeLanguage, _ := k.messageService.GetMessage(ctx, userEntity, "btn choose language")
	buttonSupport, _ := k.messageService.GetMessage(ctx, userEntity, "btn support")

	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonCreateAppointment, "create_appointment"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonMyAppointments, "my_appointments"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonChangeLanguage, "change_language"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(buttonSupport, "https://t.me/Unitedmedclinic"),
		),
	)

	msgText, _ = k.messageService.GetMessage(ctx, userEntity, "Start")
	return msgText, keyboard
}

func (k Keyboards) ConfigureDoctorInfoMessage(ctx context.Context, userEntity entity.User, doctorId int) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	doctor, err := k.appointmentService.GetDoctorInfo(ctx, userEntity, doctorId)
	msgText, _ = k.messageService.GetMessage(ctx, userEntity, "empty doctor info")

	buttonBackText, _ := k.messageService.GetMessage(ctx, userEntity, "btn back")
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonBackText, fmt.Sprintf("back_%s_%d", *userEntity.GetState(), doctorId)),
		),
	)

	if err != nil {
		return msgText, keyboard
	}

	msgText, _ = k.messageService.GetMessage(ctx, userEntity, "doctor info")
	return fmt.Sprintf(msgText, doctor.GetName(), doctor.GetInfo()), keyboard
}

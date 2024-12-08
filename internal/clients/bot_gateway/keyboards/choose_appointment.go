package keyboards

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	entity "sorkin_bot/internal/domain/entity/user"
)

func (k Keyboards) ConfigureChooseAppointmentMessage(ctx context.Context, userEntity entity.User) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	buttonClinicAppointment, err := k.messageService.GetMessage(ctx, userEntity, "btn clinic appointment")
	if err != nil {
		//	todo
	}
	buttonOnlineAppointment, err := k.messageService.GetMessage(ctx, userEntity, "btn online appointment")
	if err != nil {
		//	todo
	}
	buttonHomeAppointment, err := k.messageService.GetMessage(ctx, userEntity, "btn home appointment")
	if err != nil {
		//	todo
	}

	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonClinicAppointment, "clinic_appointment"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonOnlineAppointment, "online_appointment"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonHomeAppointment, "home_appointment"),
		),
	)

	msgText, err = k.messageService.GetMessage(ctx, userEntity, "choose appointment")
	if err != nil {
		//	todo
	}
	return msgText, keyboard
}

func (k Keyboards) ConfigureDoctorOrReasonMessage(ctx context.Context, userEntity entity.User) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	buttonChooseDoctor, err := k.messageService.GetMessage(ctx, userEntity, "btn by_doctor appointment")
	if err != nil {
		//	todo
	}
	buttonChooseReason, err := k.messageService.GetMessage(ctx, userEntity, "btn by_reason appointment")
	if err != nil {
		//	todo
	}

	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonChooseDoctor, "by_doctor"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonChooseReason, "by_reason"),
		),
	)

	msgText, err = k.messageService.GetMessage(ctx, userEntity, "doctor or reason appointment")
	if err != nil {
		//	todo
	}
	return msgText, keyboard
}

func (k Keyboards) ConfigureChooseHomeDoctorSpecificationMessage(ctx context.Context, userEntity entity.User) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	buttonChooseDoctor, err := k.messageService.GetMessage(ctx, userEntity, "btn pediatrician specification")
	if err != nil {
		//	todo
	}
	buttonChooseReason, err := k.messageService.GetMessage(ctx, userEntity, "btn therapist specification")
	if err != nil {
		//	todo
	}

	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonChooseDoctor, "pediatrician"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonChooseReason, "therapist"),
		),
	)

	msgText, err = k.messageService.GetMessage(ctx, userEntity, "choose home doctor specification")
	if err != nil {
		//	todo
	}
	return msgText, keyboard
}

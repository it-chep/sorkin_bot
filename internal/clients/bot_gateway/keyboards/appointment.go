package keyboards

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"strings"
)

func (k Keyboards) ConfigureFastAppointmentMessage(
	ctx context.Context,
	userEntity entity.User,
	schedulesMap map[int]appointment.Schedule,
) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msgText, err := k.messageService.GetMessage(ctx, userEntity, "Choose fast appointment")
	if err != nil {
		return msgText, keyboard
	}

	translatedSpecialities, _ := k.messageService.GetTranslationsBySlugKeyProfession(ctx, "Дополнительно:")

	for doctorId, schedule := range schedulesMap {
		for _, professionSlug := range strings.Split(schedule.GetProfession(), ",") {
			trimmedProfession := strings.TrimSpace(professionSlug)
			if speciality, ok := translatedSpecialities[trimmedProfession]; ok {
				btn := tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s || %s", schedule.GetTimeStartShort(), schedule.GetDoctorName()),
					fmt.Sprintf("fast__%d__%d__%s__%s", *speciality.GetSourceId(), doctorId, schedule.GetTimeStart(), schedule.GetTimeEnd()),
				)
				row := tgbotapi.NewInlineKeyboardRow(btn)
				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
				break
			}
		}
	}

	return msgText, keyboard
}

func (k Keyboards) ConfigureGetMyAppointmentsMessage(
	ctx context.Context,
	userEntity entity.User,
	appointments []appointment.Appointment,
	offset int,
) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msgText, err := k.messageService.GetMessage(ctx, userEntity, "Select appointment")

	if err != nil {
		return msgText, keyboard
	}

	lengthOfAppointments := len(appointments[offset:])
	if lengthOfAppointments > InlineButtonsLimit {
		count := 0
		for _, appointmentEntity := range appointments[offset:] {
			if count == InlineButtonsLimit {
				break
			}
			btn := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s: %s - %s", appointmentEntity.GetDate(), appointmentEntity.GetTimeStartShort(), appointmentEntity.GetTimeEndShort()),
				fmt.Sprintf("appointmentId_%d", appointmentEntity.Id()),
			)
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			count++
		}
	} else {
		for _, appointmentEntity := range appointments[offset:] {
			btn := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s: %s - %s", appointmentEntity.GetDate(), appointmentEntity.GetTimeStartShort(), appointmentEntity.GetTimeEndShort()),
				fmt.Sprintf("appointmentId_%d", appointmentEntity.Id()),
			)
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
	}

	keyboard = k.moreLessButtons(offset, lengthOfAppointments, keyboard)

	return msgText, keyboard
}

func (k Keyboards) ConfigureConfirmAppointmentMessage(ctx context.Context, userEntity entity.User, doctorId int) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	buttonTextYes, _ := k.messageService.GetMessage(ctx, userEntity, "confirm appointment ? btn yes")
	buttonTextNo, _ := k.messageService.GetMessage(ctx, userEntity, "confirm appointment ? btn no")
	buttonDoc, _ := k.messageService.GetMessage(ctx, userEntity, "doc information button")
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonDoc, fmt.Sprintf("doc_info_%d", doctorId)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonTextYes, "confirm_appointment"),
			tgbotapi.NewInlineKeyboardButtonData(buttonTextNo, "reject_appointment"),
		),
	)
	msgText, _ = k.messageService.GetMessage(ctx, userEntity, "confirm appointment ? text")
	return msgText, keyboard
}

func (k Keyboards) ConfigureAppointmentDetailMessage(ctx context.Context, userEntity entity.User, appointmentEntity appointment.Appointment) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	var cancelText, docText, exitText, rescheduleText string
	var err error

	unformattedText, _ := k.messageService.GetMessage(ctx, userEntity, "detail appointment")
	msgText = fmt.Sprintf(unformattedText, appointmentEntity.GetStringDateTimeStart())

	cancelText, err = k.messageService.GetMessage(ctx, userEntity, "cancel appointment button")
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}
	}
	docText, err = k.messageService.GetMessage(ctx, userEntity, "doc information button")
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}
	}
	exitText, err = k.messageService.GetMessage(ctx, userEntity, "btn exit")
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}
	}
	rescheduleText, err = k.messageService.GetMessage(ctx, userEntity, "reschedule appointment button")
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}
	}

	// формируем клавиатуру действий с онлайн записью
	cancelAppointmentButton := tgbotapi.NewInlineKeyboardButtonData(cancelText, fmt.Sprintf("cancel_%d", appointmentEntity.Id()))
	rescheduleAppointmentButton := tgbotapi.NewInlineKeyboardButtonData(rescheduleText, fmt.Sprintf("reschedule_%d", appointmentEntity.Id()))

	docBtn := tgbotapi.NewInlineKeyboardButtonData(docText, fmt.Sprintf("doc_info_%d", appointmentEntity.DoctorId()))
	exitBtn := tgbotapi.NewInlineKeyboardButtonData(exitText, "exit")
	keyboardRowDoctor := tgbotapi.NewInlineKeyboardRow(docBtn)
	keyboardRowCancel := tgbotapi.NewInlineKeyboardRow(cancelAppointmentButton, rescheduleAppointmentButton)
	keyboardRowExit := tgbotapi.NewInlineKeyboardRow(exitBtn)

	keyboard = tgbotapi.NewInlineKeyboardMarkup(keyboardRowDoctor, keyboardRowCancel, keyboardRowExit)

	return msgText, keyboard
}

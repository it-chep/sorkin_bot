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

	translatedSpecialities, _ := k.messageService.GetTranslationsBySlugKeyProfession(ctx, "Врач")

	for doctorId, schedule := range schedulesMap {
		for _, professionSlug := range strings.Split(schedule.GetProfession(), ",") {
			trimmedProfession := strings.TrimSpace(professionSlug)
			if speciality, ok := translatedSpecialities[trimmedProfession]; ok {
				langCode := *userEntity.GetLanguageCode()
				translatedSpeciality := k.appointmentService.GetSpecialityTranslate(langCode, speciality)

				btn := tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s || %s || %s", schedule.GetTimeStartShort(), translatedSpeciality, schedule.GetDoctorName()),
					fmt.Sprintf("fast__%d__%s__%s", doctorId, schedule.GetTimeStart(), schedule.GetTimeEnd()),
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
	msgText, err := k.messageService.GetMessage(ctx, userEntity, "Choose speciality")

	if err != nil {
		return msgText, keyboard
	}

	lengthOfAppointments := len(appointments)
	if lengthOfAppointments > InlineButtonsLimit {
		count := 0
		for _, appointmentEntity := range appointments {
			if count == InlineButtonsLimit {
				break
			}
			btn := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s: %s - %s", appointmentEntity.GetDate(), appointmentEntity.GetTimeStartShort(), appointmentEntity.GetTimeEndShort()),
				fmt.Sprintf("appointmentId_%d", appointmentEntity.GetAppointmentId()),
			)
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			count++
		}
	} else {
		for _, appointmentEntity := range appointments {
			btn := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s: %s - %s", appointmentEntity.GetDate(), appointmentEntity.GetTimeStartShort(), appointmentEntity.GetTimeEndShort()),
				fmt.Sprintf("appointmentId_%d", appointmentEntity.GetAppointmentId()),
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

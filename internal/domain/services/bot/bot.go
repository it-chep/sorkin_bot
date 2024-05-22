package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"strings"
)

const (
	InlineButtonsLimit = 10
)

type BotService struct {
	logger                 *slog.Logger
	messageService         messageService
	appointmentService     appointmentService
	readTranslationStorage readTranslationStorage
}

func NewBotService(logger *slog.Logger, messageService messageService, appointmentService appointmentService, readTranslationStorage readTranslationStorage) BotService {
	return BotService{
		logger:                 logger,
		messageService:         messageService,
		appointmentService:     appointmentService,
		readTranslationStorage: readTranslationStorage,
	}
}

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

func (bs BotService) ConfigureConfirmAppointmentMessage(ctx context.Context, userEntity entity.User, doctorId int) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	buttonTextYes, _ := bs.messageService.GetMessage(ctx, userEntity, "confirm appointment ? btn yes")
	buttonTextNo, _ := bs.messageService.GetMessage(ctx, userEntity, "confirm appointment ? btn no")
	buttonDoc, _ := bs.messageService.GetMessage(ctx, userEntity, "doc information button")
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonDoc, fmt.Sprintf("doc_info_%d", doctorId)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonTextYes, "confirm_appointment"),
			tgbotapi.NewInlineKeyboardButtonData(buttonTextNo, "reject_appointment"),
		),
	)
	msgText, _ = bs.messageService.GetMessage(ctx, userEntity, "confirm appointment ? text")
	return msgText, keyboard
}

func (bs BotService) ConfigureGetSpecialityMessage(
	ctx context.Context,
	userEntity entity.User,
	translatedSpecialities map[int]string,
	offset int,
) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msgText, err := bs.messageService.GetMessage(ctx, userEntity, "Choose speciality")
	if err != nil {
		return msgText, keyboard
	}
	lengthOfSpecialities := len(translatedSpecialities)
	if lengthOfSpecialities > InlineButtonsLimit {
		count := 0
		for specialityId, translatedSpeciality := range translatedSpecialities {
			if count == InlineButtonsLimit {
				break
			}
			btn := tgbotapi.NewInlineKeyboardButtonData(translatedSpeciality, fmt.Sprintf("%d", specialityId))
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			count++
		}
	} else {
		for specialityId, translatedSpeciality := range translatedSpecialities {
			btn := tgbotapi.NewInlineKeyboardButtonData(translatedSpeciality, fmt.Sprintf("%d", specialityId))
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
	}

	keyboard = bs.moreLessButtons(offset, lengthOfSpecialities, keyboard)

	return msgText, keyboard
}

func (bs BotService) ConfigureGetDoctorMessage(
	ctx context.Context,
	userEntity entity.User,
	doctors map[int]string,
	offset int,
) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msgText, err := bs.messageService.GetMessage(ctx, userEntity, "Choose doctor")
	if err != nil {
		return msgText, keyboard
	}
	lengthOfDoctors := len(doctors)
	if lengthOfDoctors > InlineButtonsLimit {
		count := 0
		for doctorId, doctorName := range doctors {
			if count == InlineButtonsLimit {
				break
			}
			btn := tgbotapi.NewInlineKeyboardButtonData(doctorName, fmt.Sprintf("%d", doctorId))
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			count++
		}
	} else if lengthOfDoctors == 0 {
		msgText, err = bs.messageService.GetMessage(ctx, userEntity, "empty doctors")
		if err != nil {
			return msgText, keyboard
		}
	} else {
		fmt.Println(fmt.Sprintf("doctors: %v", doctors))
		for doctorId, doctorName := range doctors {
			btn := tgbotapi.NewInlineKeyboardButtonData(doctorName, fmt.Sprintf("%d", doctorId))
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
	}

	keyboard = bs.moreLessButtons(offset, lengthOfDoctors, keyboard)

	return msgText, keyboard
}

func (bs BotService) ConfigureGetScheduleMessage(
	ctx context.Context,
	userEntity entity.User,
	schedules []appointment.Schedule,
	offset int,
) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msgText, err := bs.messageService.GetMessage(ctx, userEntity, "Choose schedule")
	if err != nil {
		return msgText, keyboard
	}
	lengthOfSchedules := len(schedules)

	if lengthOfSchedules > InlineButtonsLimit {
		count := 0
		for _, schedule := range schedules[offset:] {
			if count == InlineButtonsLimit {
				break
			}
			btn := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s: %s - %s", schedule.GetDate(), schedule.GetTimeStartShort(), schedule.GetTimeEndShort()),
				fmt.Sprintf("schedule_%d_%s_%s_%s", schedule.GetDoctorId(), schedule.GetTimeStart(), schedule.GetTimeEnd(), schedule.GetDate()))
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			count++
		}
	} else {
		for _, schedule := range schedules {
			btn := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s: %s - %s", schedule.GetDate(), schedule.GetTimeStartShort(), schedule.GetTimeEndShort()),
				fmt.Sprintf("schedule_%d_%s_%s_%s", schedule.GetDoctorId(), schedule.GetTimeStart(), schedule.GetTimeEnd(), schedule.GetDate()))
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
	}

	keyboard = bs.moreLessButtons(offset, lengthOfSchedules, keyboard)

	return msgText, keyboard
}

func (bs BotService) ConfigureGetMyAppointmentsMessage(
	ctx context.Context,
	userEntity entity.User,
	appointments map[int]string,
	offset int,
) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msgText, err := bs.messageService.GetMessage(ctx, userEntity, "Choose speciality")
	if err != nil {
		return msgText, keyboard
	}
	lengthOfAppointments := len(appointments)
	if lengthOfAppointments > InlineButtonsLimit {
		count := 0
		for appointmentId, appointmentText := range appointments {
			if count == InlineButtonsLimit {
				break
			}
			btn := tgbotapi.NewInlineKeyboardButtonData(appointmentText, fmt.Sprintf("%d", appointmentId))
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			count++
		}
	} else {
		for appointmentId, appointmentText := range appointments {
			btn := tgbotapi.NewInlineKeyboardButtonData(appointmentText, fmt.Sprintf("%d", appointmentId))
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
	}

	keyboard = bs.moreLessButtons(offset, lengthOfAppointments, keyboard)

	return msgText, keyboard
}

func (bs BotService) moreLessButtons(offset, lengthOfItems int, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.InlineKeyboardMarkup {

	// Добавляем кнопки с переключением
	if offset > lengthOfItems {
		btnBack := tgbotapi.NewInlineKeyboardButtonData("<", fmt.Sprintf("offset_%d_<", offset))
		row := tgbotapi.NewInlineKeyboardRow(btnBack)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	} else if offset == 0 && lengthOfItems > InlineButtonsLimit {
		btnMore := tgbotapi.NewInlineKeyboardButtonData(">", fmt.Sprintf("offset_%d_>", offset))
		row := tgbotapi.NewInlineKeyboardRow(btnMore)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	} else if offset != 0 && lengthOfItems > InlineButtonsLimit {
		btnBack := tgbotapi.NewInlineKeyboardButtonData("<", fmt.Sprintf("offset_%d_<", offset))
		btnMore := tgbotapi.NewInlineKeyboardButtonData(">", fmt.Sprintf("offset_%d_>", offset))
		row := tgbotapi.NewInlineKeyboardRow(btnBack, btnMore)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	return keyboard
}

func (bs BotService) ConfigureFastAppointmentMessage(
	ctx context.Context,
	userEntity entity.User,
	schedulesMap map[int]appointment.Schedule,
) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msgText, err := bs.messageService.GetMessage(ctx, userEntity, "Choose fast appointment")
	if err != nil {
		return msgText, keyboard
	}

	translatedSpecialities, _ := bs.readTranslationStorage.GetTranslationsBySlugKeyProfession(ctx, "Врач")

	for doctorId, schedule := range schedulesMap {
		for _, professionSlug := range strings.Split(schedule.GetProfession(), ",") {
			trimmedProfession := strings.TrimSpace(professionSlug)
			if speciality, ok := translatedSpecialities[trimmedProfession]; ok {
				langCode := *userEntity.GetLanguageCode()
				translatedSpeciality := bs.appointmentService.GetSpecialityTranslate(langCode, speciality)

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

package keyboards

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"sort"
)

const (
	InlineButtonsLimit = 10
)

type Keyboards struct {
	logger             *slog.Logger
	messageService     messageService
	appointmentService appointmentService
}

func NewKeyboards(logger *slog.Logger, messageService messageService, appointmentService appointmentService) Keyboards {
	return Keyboards{
		logger:             logger,
		messageService:     messageService,
		appointmentService: appointmentService,
	}
}

func (k Keyboards) ConfigureGetSpecialityMessage(
	ctx context.Context,
	userEntity entity.User,
	translatedSpecialities map[int]string,
	offset int,
) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msgText, err := k.messageService.GetMessage(ctx, userEntity, "Choose speciality")
	if err != nil {
		return msgText, keyboard
	}

	keys := make([]int, 0, len(translatedSpecialities))
	for k := range translatedSpecialities {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	lengthOfSpecialities := len(translatedSpecialities)
	if lengthOfSpecialities > InlineButtonsLimit {
		count := 0
		for _, specialityId := range keys {
			if count == InlineButtonsLimit {
				break
			}
			translatedSpeciality := translatedSpecialities[specialityId]
			btn := tgbotapi.NewInlineKeyboardButtonData(translatedSpeciality, fmt.Sprintf("%d", specialityId))
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			count++
		}
	} else {
		for _, specialityId := range keys {
			translatedSpeciality := translatedSpecialities[specialityId]
			btn := tgbotapi.NewInlineKeyboardButtonData(translatedSpeciality, fmt.Sprintf("%d", specialityId))
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
	}

	keyboard = k.moreLessButtons(offset, lengthOfSpecialities, keyboard)

	return msgText, keyboard
}

func (k Keyboards) ConfigureGetDoctorMessage(
	ctx context.Context,
	userEntity entity.User,
	doctors map[int]string,
	offset int,
) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msgText, err := k.messageService.GetMessage(ctx, userEntity, "Choose doctor")
	if err != nil {
		return msgText, keyboard
	}

	keys := make([]int, 0, len(doctors))
	for k := range doctors {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	lengthOfDoctors := len(doctors)
	if lengthOfDoctors > InlineButtonsLimit {
		count := 0
		for _, doctorId := range keys {
			if count == InlineButtonsLimit {
				break
			}
			doctorName := doctors[doctorId]
			btn := tgbotapi.NewInlineKeyboardButtonData(doctorName, fmt.Sprintf("%d_%s", doctorId, doctorName))
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			count++
		}
	} else {
		for _, doctorId := range keys {
			doctorName := doctors[doctorId]
			btn := tgbotapi.NewInlineKeyboardButtonData(doctorName, fmt.Sprintf("%d_%s", doctorId, doctorName))
			row := tgbotapi.NewInlineKeyboardRow(btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
	}

	keyboard = k.moreLessButtons(offset, lengthOfDoctors, keyboard)

	return msgText, keyboard
}

func (k Keyboards) ConfigureGetScheduleMessage(
	ctx context.Context,
	userEntity entity.User,
	schedules []appointment.Schedule,
	offset int,
) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msgText, err := k.messageService.GetMessage(ctx, userEntity, "Choose schedule")
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

	keyboard = k.moreLessButtons(offset, lengthOfSchedules, keyboard)

	return msgText, keyboard
}

func (k Keyboards) moreLessButtons(offset, lengthOfItems int, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.InlineKeyboardMarkup {

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

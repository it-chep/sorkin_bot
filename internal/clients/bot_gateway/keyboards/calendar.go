package keyboards

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	entity "sorkin_bot/internal/domain/entity/user"
	"time"
)

var dayIndex = map[time.Weekday]int{
	time.Monday:    0,
	time.Tuesday:   1,
	time.Wednesday: 2,
	time.Thursday:  3,
	time.Friday:    4,
	time.Saturday:  0,
	time.Sunday:    6,
}

func (k Keyboards) GenerateCalendarKeyboard(ctx context.Context, user entity.User, year int, month time.Month, schedulesMap map[time.Time]bool) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	var (
		navigationRow []tgbotapi.InlineKeyboardButton
	)

	now := time.Now()
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

	row := make([]tgbotapi.InlineKeyboardButton, 0, 7)
	weekRow := make([]tgbotapi.InlineKeyboardButton, 0, 7)
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 8)

	// Получаем переводы на клавиатуру
	daysOfWeek, err := k.messageService.GetWeekdaysName(ctx, user)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}
	}
	// Создаем строку из ПН ВТ СР ЧТ ПТ СБ ВС
	for _, day := range daysOfWeek {
		weekRow = append(weekRow, tgbotapi.NewInlineKeyboardButtonData(day, "ignore"))
	}
	rows = append(rows, weekRow)

	// Заполнение пустыми значениями дни до 1 числа месяца
	for i := 0; i < dayIndex[startDate.Weekday()]; i++ {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", "ignore"))
	}
	emptyButton := 0
	// Проход по дням
	for day := 1; day <= daysInMonth; day++ {
		// Преобразуем день в callback_data необходимую для нашей логики
		dateStr := fmt.Sprintf("%d-%02d-%02d", year, month, day)

		// Дата в тип даты а не строки
		parsedDate, _ := time.Parse(time.DateOnly, dateStr)
		weekday := parsedDate.Weekday()

		if startDate.Weekday() == time.Saturday && (day == 1 || day == 2) {
			continue
		}

		if ok, _ := schedulesMap[parsedDate]; !ok && len(schedulesMap) > 0 {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", "ignore"))
			emptyButton++
		} else if weekday == time.Saturday || weekday == time.Sunday {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", "ignore"))
			emptyButton++
		} else if now.Before(parsedDate) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", day), dateStr))
		} else {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", "ignore"))
			emptyButton++
		}

		// Если уже заполнены все дни недели, то добавляем новую строку
		if len(row) == 7 {
			emptyButton = 0
			emptyDay := 0
			for _, button := range row {
				if button.Text == " " {
					emptyDay++
				}
			}
			if emptyDay == 7 {
				row = make([]tgbotapi.InlineKeyboardButton, 0, 7)
				continue
			}
			rows = append(rows, row)
			row = make([]tgbotapi.InlineKeyboardButton, 0, 7)
		}
	}

	//Если в последней строке все кнопки пустые то не стоит ее дозаполнять
	if len(row) != emptyButton {
		// Если последний день месяца - не ВС то добавляем пустые ячейки
		for len(row) < 7 {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", "ignore"))
		}

		rows = append(rows, row)
	}

	prevMonth := month - 1
	nextMonth := month + 1
	// Получаем название месяца в языке
	monthName, err := k.messageService.GetMessage(ctx, user, month.String())
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}
	}
	nowMonth := time.Now().Month()

	// Создаем клавиатуру в зависимости от месяца. У нас должна быть возможность только на текущий и следущий месяц записаться
	if nowMonth == month {
		navigationRow = []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(" ", "ignore"),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %d", monthName, year), "ignore"),
			tgbotapi.NewInlineKeyboardButtonData("→", fmt.Sprintf("calendar:%d-%02d", year, nextMonth)),
		}
	} else {
		navigationRow = []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("←", fmt.Sprintf("calendar:%d-%02d", year, prevMonth)),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %d", monthName, year), "ignore"),
			tgbotapi.NewInlineKeyboardButtonData(" ", "ignore"),
		}
	}

	rows = append(rows, navigationRow)

	msgText, err = k.messageService.GetMessage(ctx, user, "choose date")
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}
	}
	return msgText, tgbotapi.NewInlineKeyboardMarkup(rows...)
}

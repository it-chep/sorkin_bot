package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"strconv"
	"strings"
	"time"
)

func (c *CallbackBotMessage) getCalendar(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	var (
		year         int
		month        time.Month
		err          error
		schedulesMap map[time.Time]bool
	)
	callbackDataItems := strings.Split(callbackData, "-")

	year, err = strconv.Atoi(callbackDataItems[0])
	if err != nil {
		year = time.Now().Year()
	}

	monthInt, err := strconv.Atoi(callbackDataItems[1])
	if err != nil || monthInt < 1 {
		month = time.Now().Month()
	} else if monthInt > 12 {
		month = time.January
		year = time.Now().Year() + 1
	} else {
		month = time.Month(monthInt)
	}

	draftAppointment, err := c.appointmentService.GetDraftAppointment(ctx, userEntity.GetTgId())
	if err != nil {
		return
	}
	schedulesMap = make(map[time.Time]bool)
	if doctorId := draftAppointment.GetDoctorId(); doctorId != nil {
		parsedDate, _ := time.Parse(time.DateOnly, fmt.Sprintf("%d-%02d-01", year, month))
		schedulesMap, err = c.appointmentService.GetSchedulePeriodsByDoctorId(ctx, *doctorId, parsedDate)
		if len(schedulesMap) == 0 {
			c.botGateway.SendEmptySchedulePeriods(ctx, userEntity, messageDTO)
			c.getDefaultCalendar(ctx, messageDTO, userEntity)
			return
		}
	}
	c.botGateway.SendCalendarMessage(ctx, userEntity, messageDTO, year, month, schedulesMap)
}

func (c *CallbackBotMessage) getDefaultCalendar(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User) {
	nowYear, nowMonth := time.Now().Year(), time.Now().Month()
	draftAppointment, err := c.appointmentService.GetDraftAppointment(ctx, userEntity.GetTgId())
	if err != nil {
		return
	}
	schedulesMap := make(map[time.Time]bool)
	if doctorId := draftAppointment.GetDoctorId(); doctorId != nil {
		schedulesMap, err = c.appointmentService.GetSchedulePeriodsByDoctorId(ctx, *doctorId, time.Now())
		if len(schedulesMap) == 0 {
			c.botGateway.SendEmptySchedulePeriods(ctx, userEntity, messageDTO)
			return
		}
	}
	c.botGateway.SendCalendarMessage(ctx, userEntity, messageDTO, nowYear, nowMonth, schedulesMap)
}

func (c *CallbackBotMessage) convertCallbackDateToDate(callbackData string) (date time.Time, err error) {
	callbackDataItems := strings.Split(callbackData, "-")

	if len(callbackDataItems) != 3 {
		return time.Time{}, fmt.Errorf("invalid callback data format")
	}
	year, err := strconv.Atoi(callbackDataItems[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year format: %v", err)
	}
	month, err := strconv.Atoi(callbackDataItems[1])
	if err != nil || month < 1 || month > 12 {
		return time.Time{}, fmt.Errorf("invalid month format: %v", err)
	}

	day, err := strconv.Atoi(callbackDataItems[2])
	if err != nil || day < 1 || day > 31 {
		return time.Time{}, fmt.Errorf("invalid day format: %v", err)
	}

	date = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return date, nil
}

func (c *CallbackBotMessage) getSchedules(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	var msgText string
	var err error
	var msg tgbotapi.MessageConfig

	now := time.Now()
	callbackDate, err := c.convertCallbackDateToDate(callbackData)

	if callbackDate.Before(now) {
		schedulesMap := make(map[time.Time]bool)
		c.botGateway.SendForbiddenAction(ctx, userEntity, messageDTO)
		c.botGateway.SendCalendarMessage(ctx, userEntity, messageDTO, now.Year(), now.Month(), schedulesMap)
		return
	}

	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))

	msgText, err = c.messageService.GetMessage(ctx, userEntity, "your doctor")
	draft, err := c.appointmentService.GetDraftAppointment(ctx, userEntity.GetTgId())
	if err != nil {
		// todo
		return
	}
	msg = tgbotapi.NewMessage(c.tgUser.TgID, fmt.Sprintf(msgText, *draft.GetDoctorName()))
	c.bot.SendMessage(msg, messageDTO)

	sentMessageId := c.botGateway.SendWaitMessage(ctx, userEntity, messageDTO, "wait schedules")

	schedules, err := c.appointmentService.GetSchedulesByDoctorId(ctx, userEntity, callbackDate, draft.GetDoctorId())

	if err != nil {
		msgText, err = c.messageService.GetMessage(ctx, userEntity, "empty schedules")
		msg = tgbotapi.NewMessage(c.tgUser.TgID, msgText)
		c.bot.SendMessage(msg, messageDTO)
		c.machine.SetState(userEntity, state_machine.ChooseDoctor)

		draftAppointmentEntity, _ := c.appointmentService.GetDraftAppointment(ctx, userEntity.GetTgId())
		if draftAppointmentEntity.GetSpecialityId() == nil {
			c.getDoctors(ctx, messageDTO, userEntity)
			return
		}
		c.getDoctorsBySpecialityId(ctx, messageDTO, userEntity, *draftAppointmentEntity.GetSpecialityId())
		return
	}

	c.bot.RemoveMessage(c.tgUser.TgID, sentMessageId)
	c.botGateway.SendSchedulesMessage(ctx, userEntity, messageDTO, schedules, ZeroOffset)
	c.machine.SetState(userEntity, state_machine.ChooseSchedule)
}

func (c *CallbackBotMessage) chooseSchedules(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	if strings.Contains(callbackData, "offset") {
		c.moreLessSchedules(ctx, messageDTO, userEntity, callbackData)
	} else if strings.Contains(callbackData, "schedule") {
		c.saveDraftAppointment(ctx, messageDTO, userEntity, callbackData)
	} else {
		c.getSchedules(ctx, messageDTO, userEntity, callbackData)
	}
}

func (c *CallbackBotMessage) chooseCalendar(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
	callbackDataItems := strings.Split(callbackData, "-")
	if len(callbackDataItems) == 2 {
		c.getCalendar(ctx, messageDTO, userEntity, callbackData)
		c.machine.SetState(userEntity, state_machine.ChooseCalendar)
	} else if len(callbackDataItems) == 3 {
		c.getSchedules(ctx, messageDTO, userEntity, callbackData)
	} else if callbackData == "ignore" {
		c.getDefaultCalendar(ctx, messageDTO, userEntity)
		c.machine.SetState(userEntity, state_machine.ChooseCalendar)
	} else {
		c.getDefaultCalendar(ctx, messageDTO, userEntity)
		c.machine.SetState(userEntity, state_machine.ChooseCalendar)
	}
}

func (c *CallbackBotMessage) moreLessSchedules(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	offset, _ := strconv.Atoi(strings.Split(callbackData, "_")[1])
	if strings.Contains(callbackData, ">") {
		offset += 10
	} else {
		offset -= 10
	}

	schedules, err := c.appointmentService.GetSchedulesByDoctorId(ctx, userEntity, time.Now(), nil)
	if err != nil {
		return
	}

	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
	c.botGateway.SendSchedulesMessage(ctx, userEntity, messageDTO, schedules, offset)
}

func (c *CallbackBotMessage) saveDraftAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	scheduleItems := strings.Split(callbackData, "_")
	doctorId, _ := strconv.Atoi(scheduleItems[1])
	fullTimeStart := scheduleItems[2]
	fullTimeEnd := scheduleItems[3]
	date := scheduleItems[4]

	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
	// todo тут какой-то кринж потому что сначала update а потом select идет, можно все в 1 методе
	c.appointmentService.UpdateDraftAppointmentDate(ctx, userEntity.GetTgId(), fullTimeStart, fullTimeEnd, date)
	draftAppointment, err := c.appointmentService.GetDraftAppointment(ctx, userEntity.GetTgId())
	if err != nil {
		// todo
		return
	}

	if *draftAppointment.GetAppointmentType() == appointment.HomeAppointment {
		c.appointmentService.UpdateDraftAppointmentIntField(ctx, userEntity.GetTgId(), doctorId, "doctor_id")
		c.appointmentService.UpdateDraftAppointmentDoctorName(ctx, userEntity.GetTgId(), doctorId)
	}

	if userEntity.GetPhone() == nil {
		c.botGateway.SendGetPhoneMessage(ctx, userEntity, messageDTO)
		c.machine.SetState(userEntity, state_machine.GetPhone)
		return
	}

	if *draftAppointment.GetAppointmentType() == appointment.HomeAppointment {
		c.botGateway.SendGetHomeAddressMessage(ctx, userEntity, messageDTO)
		c.machine.SetState(userEntity, state_machine.SetAddress)
		return
	}

	c.botGateway.SendConfirmAppointmentMessage(ctx, userEntity, messageDTO, doctorId)
	c.machine.SetState(userEntity, state_machine.CreateAppointment)
}

package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/config"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"strconv"
	"strings"
	"time"
)

func (c *CallbackBotMessage) getClinicConfig() config.MISConfig {
	return config.NewConfig().MIS
}

func (c *CallbackBotMessage) preCreateAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	if callbackData == "confirm_appointment" {
		c.confirmAppointment(ctx, messageDTO, userEntity)
	} else if callbackData == "reject_appointment" {
		c.appointmentService.CleanDraftAppointment(ctx, userEntity.GetTgId())
		c.rejectAppointment(ctx, messageDTO, userEntity)
	} else if strings.Contains(callbackData, "doc_info") {
		c.getDoctorInfo(ctx, messageDTO, userEntity, callbackData)
	}
}

func (c *CallbackBotMessage) rejectAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User) {
	msgText, err := c.messageService.GetMessage(ctx, userEntity, "have rejected appointment")
	if err != nil {
		return
	}
	msg := tgbotapi.NewMessage(userEntity.GetTgId(), msgText)
	c.bot.SendMessage(msg, messageDTO)
	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))
	c.botGateway.SendStartMessage(ctx, userEntity, messageDTO)
	c.machine.SetState(userEntity, state_machine.Start)
}

func (c *CallbackBotMessage) confirmAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User) {
	misConfig := c.getClinicConfig()
	draftAppointmentEntity, err := c.appointmentService.GetDraftAppointment(ctx, userEntity.GetTgId())
	if err != nil {
		return
	}
	if draftAppointmentEntity.GetDoctorId() == nil {
		return
	}

	//todo херня, надо улучшить
	appointmentString := fmt.Sprintf("doctorId_%d__timeStart_%s__timeEnd_%s",
		*draftAppointmentEntity.GetDoctorId(),
		*draftAppointmentEntity.GetTimeStart(),
		*draftAppointmentEntity.GetTimeEnd(),
	)

	appointmentId := c.appointmentService.CreateAppointment(ctx, userEntity, draftAppointmentEntity, appointmentString)
	if appointmentId != nil {
		c.appointmentService.UpdateDraftAppointmentStatus(ctx, userEntity.GetTgId(), *appointmentId)
	}

	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))

	switch *draftAppointmentEntity.GetAppointmentType() {
	case appointment.ClinicAppointment:
		msgText, _ := c.messageService.GetMessage(ctx, userEntity, "successfully created appointment")
		msg := tgbotapi.NewMessage(c.tgUser.TgID, fmt.Sprintf(msgText, *draftAppointmentEntity.GetDoctorName(), *draftAppointmentEntity.GetTimeStart()))
		c.bot.SendMessage(msg, messageDTO)
		c.bot.SendLocation(userEntity.GetTgId(), misConfig.Latitude, misConfig.Longitude, messageDTO)
	case appointment.HomeAppointment:
		msgText, _ := c.messageService.GetMessage(ctx, userEntity, "successfully created home appointment")
		msg := tgbotapi.NewMessage(c.tgUser.TgID, fmt.Sprintf(msgText, *draftAppointmentEntity.GetDoctorName(), *draftAppointmentEntity.GetTimeStart()))
		c.bot.SendMessage(msg, messageDTO)
	case appointment.OnlineAppointment:
		msgText, _ := c.messageService.GetMessage(ctx, userEntity, "successfully created online appointment")
		msg := tgbotapi.NewMessage(c.tgUser.TgID, fmt.Sprintf(msgText, *draftAppointmentEntity.GetDoctorName(), *draftAppointmentEntity.GetTimeStart()))
		c.bot.SendMessage(msg, messageDTO)
	}

	c.botGateway.SendStartMessage(ctx, userEntity, messageDTO)
	c.machine.SetState(userEntity, state_machine.Start)
}

func (c *CallbackBotMessage) fastAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	if strings.Contains(callbackData, "fast_") {
		c.bot.RemoveMessage(userEntity.GetTgId(), int(messageDTO.MessageID))

		items := strings.Split(callbackData, "__")
		specialityId, err := strconv.Atoi(items[1])
		if err != nil {
			return
		}
		doctorId, err := strconv.Atoi(items[2])
		if err != nil {
			return
		}
		timeStart := items[3]
		timeEnd := items[4]
		c.appointmentService.FastUpdateDraftAppointment(ctx, userEntity.GetTgId(), specialityId, doctorId, timeStart, timeEnd)

		if userEntity.GetPhone() != nil {
			c.botGateway.SendConfirmAppointmentMessage(ctx, userEntity, messageDTO, doctorId)
			c.machine.SetState(userEntity, state_machine.CreateAppointment)
		} else {
			c.botGateway.SendGetPhoneMessage(ctx, userEntity, messageDTO)
			c.machine.SetState(userEntity, state_machine.GetPhone)
		}
	}
}

func (c *CallbackBotMessage) getDoctorInfo(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	callbackItems := strings.Split(callbackData, "_")
	doctorId, _ := strconv.Atoi(callbackItems[2])
	c.botGateway.SendDoctorInfoMessage(ctx, userEntity, messageDTO, int(messageDTO.MessageID), doctorId)
	c.machine.SetState(userEntity, state_machine.GetDoctorInfo)
}

func (c *CallbackBotMessage) chooseAppointmentVariant(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	c.botGateway.SendChooseAppointmentMessage(ctx, userEntity, messageDTO)
}

func (c *CallbackBotMessage) chooseAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	c.bot.RemoveMessage(userEntity.GetTgId(), int(messageDTO.MessageID))
	switch callbackData {
	case "clinic_appointment":
		c.botGateway.SendDoctorOrReasonMessage(ctx, userEntity, messageDTO)
		c.machine.SetState(userEntity, state_machine.ClinicAppointment)
		c.appointmentService.UpdateDraftAppointmentType(ctx, userEntity.GetTgId(), appointment.ClinicAppointment)
	case "online_appointment":
		c.botGateway.SendDoctorOrReasonMessage(ctx, userEntity, messageDTO)
		c.machine.SetState(userEntity, state_machine.OnlineAppointment)
		c.appointmentService.UpdateDraftAppointmentType(ctx, userEntity.GetTgId(), appointment.OnlineAppointment)
	case "home_appointment":
		c.botGateway.SendHomeDoctorSpecificationMessage(ctx, userEntity, messageDTO)
		c.machine.SetState(userEntity, state_machine.HomeAppointment)
		c.appointmentService.UpdateDraftAppointmentType(ctx, userEntity.GetTgId(), appointment.HomeAppointment)
	}
}

func (c *CallbackBotMessage) forkDoctorReasonAppointment(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	c.bot.RemoveMessage(userEntity.GetTgId(), int(messageDTO.MessageID))
	switch callbackData {
	case "by_doctor":
		c.getDoctors(ctx, messageDTO, userEntity)
		c.machine.SetState(userEntity, state_machine.ChooseDoctor)
	case "by_reason":
		c.moreLessSpeciality(ctx, messageDTO, userEntity, callbackData)
		c.machine.SetState(userEntity, state_machine.ChooseSpeciality)
	}
}

func (c *CallbackBotMessage) forkHomeDoctorSpecialisation(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	c.bot.RemoveMessage(userEntity.GetTgId(), int(messageDTO.MessageID))
	switch callbackData {
	case "pediatrician":
		//	Педиатр
		c.getDefaultCalendar(ctx, messageDTO, userEntity)
		c.machine.SetState(userEntity, state_machine.Pediatrician)
	case "therapist":
		//	Терапевт
		c.getDefaultCalendar(ctx, messageDTO, userEntity)
		c.machine.SetState(userEntity, state_machine.Therapist)
	}
}

func (c *CallbackBotMessage) getHomeVisitSchedules(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	c.bot.RemoveMessage(userEntity.GetTgId(), int(messageDTO.MessageID))
	var err error

	now := time.Now()
	callbackDataItems := strings.Split(callbackData, "-")
	if len(callbackDataItems) == 2 {
		c.getCalendar(ctx, messageDTO, userEntity, callbackData)
		return
	} else if callbackData == "ignore" {
		c.getDefaultCalendar(ctx, messageDTO, userEntity)
		return
	}
	schedulesMap := make(map[time.Time]bool)
	callbackDate, err := c.convertCallbackDateToDate(callbackData)
	if callbackDate.Before(now) {
		c.botGateway.SendForbiddenAction(ctx, userEntity, messageDTO)
		c.botGateway.SendCalendarMessage(ctx, userEntity, messageDTO, now.Year(), now.Month(), schedulesMap)
		return
	}

	sentMessageId := c.botGateway.SendWaitMessage(ctx, userEntity, messageDTO, "wait schedules")

	schedules, err := c.appointmentService.GetSchedulesToHomeVisit(ctx, userEntity, callbackDate)

	if err != nil || len(schedules) == 0 {
		c.botGateway.SendEmptySchedulesHomeVisit(ctx, userEntity, messageDTO)
		c.botGateway.SendCalendarMessage(ctx, userEntity, messageDTO, time.Now().Year(), time.Now().Month(), schedulesMap)
		return
	}

	c.bot.RemoveMessage(c.tgUser.TgID, sentMessageId)
	c.botGateway.SendSchedulesMessage(ctx, userEntity, messageDTO, schedules, ZeroOffset)
	c.machine.SetState(userEntity, state_machine.ChooseSchedule)
}

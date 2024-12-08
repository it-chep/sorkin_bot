package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/domain/entity/user/state_machine"
	"strconv"
	"strings"
)

func (c *CallbackBotMessage) getDoctorsBySpecialityId(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, specialityId int) {
	var msgText string
	var err error

	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))

	sentMessageId := c.botGateway.SendWaitMessage(ctx, userEntity, messageDTO, "wait doctors")

	doctors := c.appointmentService.GetDoctorsBySpecialityId(ctx, userEntity.GetTgId(), 0, &specialityId)

	c.bot.RemoveMessage(c.tgUser.TgID, sentMessageId)

	msgText, err = c.messageService.GetMessage(ctx, userEntity, "your speciality")

	if err != nil {
		c.bot.SendMessage(tgbotapi.NewMessage(userEntity.GetTgId(), msgText), messageDTO)
		return
	}

	specialityText, err := c.appointmentService.TranslateSpecialityByID(ctx, userEntity, specialityId)
	if err != nil {
		return
	}

	if specialityText != "" {
		c.bot.SendMessage(tgbotapi.NewMessage(userEntity.GetTgId(), fmt.Sprintf(msgText, specialityText)), messageDTO)
	}

	if len(doctors) != 0 {
		c.botGateway.SendGetDoctorsMessage(ctx, userEntity, messageDTO, doctors, ZeroOffset)
		c.machine.SetState(userEntity, state_machine.ChooseDoctor)
	} else {
		msg := tgbotapi.NewMessage(c.tgUser.TgID, msgText)
		msgText, err = c.messageService.GetMessage(ctx, userEntity, "empty doctors")

		c.bot.SendMessage(msg, messageDTO)
		c.moreLessSpeciality(ctx, messageDTO, userEntity, "")

		c.machine.SetState(userEntity, state_machine.ChooseSpeciality)
	}
}

func (c *CallbackBotMessage) chooseDoctor(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	if strings.Contains(callbackData, "offset") {
		c.moreLessDoctors(ctx, messageDTO, userEntity, callbackData)
	} else {
		doctorId, _ := strconv.Atoi(strings.Split(callbackData, "_")[0])
		c.appointmentService.UpdateDraftAppointmentIntField(ctx, userEntity.GetTgId(), doctorId, "doctor_id")
		c.appointmentService.UpdateDraftAppointmentDoctorName(ctx, userEntity.GetTgId(), doctorId)
		c.chooseCalendar(ctx, messageDTO, userEntity, callbackData)
	}
}

func (c *CallbackBotMessage) afterDoctorInfo(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	callbackDataItems := strings.Split(callbackData, "_")
	if strings.Contains(callbackData, "back") {
		previousState := callbackDataItems[1]
		doctorId, _ := strconv.Atoi(callbackDataItems[2])
		c.bot.RemoveMessage(userEntity.GetTgId(), int(messageDTO.MessageID))
		if previousState == state_machine.DetailMyAppointment {
			appointments := c.appointmentService.GetAppointments(ctx, userEntity)
			c.botGateway.SendMyAppointmentsMessage(ctx, userEntity, appointments, messageDTO, 0)
			c.machine.SetState(userEntity, state_machine.ChooseMyAppointment)
		} else if previousState == state_machine.CreateAppointment {
			c.botGateway.SendConfirmAppointmentMessage(ctx, userEntity, messageDTO, doctorId)
			c.machine.SetState(userEntity, state_machine.CreateAppointment)
		}
	}
}

func (c *CallbackBotMessage) moreLessDoctors(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User, callbackData string) {
	var doctors map[int]string
	offset, _ := strconv.Atoi(strings.Split(callbackData, "_")[1])
	if strings.Contains(callbackData, ">") {
		offset += 10
	} else if strings.Contains(callbackData, "<") {
		offset -= 10
	}
	appointment, err := c.appointmentService.GetDraftAppointment(ctx, userEntity.GetTgId())
	if err != nil {
		//todo
		return
	}
	specialityId := appointment.GetSpecialityId()
	if specialityId == nil {
		doctors = c.appointmentService.GetDoctors(ctx, userEntity.GetTgId(), offset)
	} else {
		doctors = c.appointmentService.GetDoctorsBySpecialityId(ctx, userEntity.GetTgId(), offset, specialityId)
	}
	c.botGateway.SendGetDoctorsMessage(ctx, userEntity, messageDTO, doctors, offset)
}

func (c *CallbackBotMessage) getDoctors(ctx context.Context, messageDTO tg.MessageDTO, userEntity entity.User) {

	c.bot.RemoveMessage(c.tgUser.TgID, int(messageDTO.MessageID))

	sentMessageId := c.botGateway.SendWaitMessage(ctx, userEntity, messageDTO, "wait doctors")

	doctors := c.appointmentService.GetDoctors(ctx, userEntity.GetTgId(), 0)

	c.bot.RemoveMessage(c.tgUser.TgID, sentMessageId)

	c.botGateway.SendGetDoctorsMessage(ctx, userEntity, messageDTO, doctors, ZeroOffset)
	c.machine.SetState(userEntity, state_machine.ChooseDoctor)
}

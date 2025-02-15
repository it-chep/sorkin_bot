package bot_gateway

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"time"
)

func (bg BotGateway) SendFastAppointmentMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	var msgText string
	var msg tgbotapi.MessageConfig

	messageId := bg.SendWaitMessage(ctx, user, messageDTO, "wait fast appointment")

	schedulesMap := bg.appointmentService.GetFastAppointmentSchedules(ctx)

	msgText, keyboard := bg.keyboard.ConfigureFastAppointmentMessage(ctx, user, schedulesMap)

	bg.bot.RemoveMessage(user.GetTgId(), messageId)

	msg = tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	bg.bot.SendMessage(msg, messageDTO)
}

func (bg BotGateway) SendConfirmAppointmentMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, doctorId int) {
	msgText, keyboard := bg.keyboard.ConfigureConfirmAppointmentMessage(ctx, user, doctorId)
	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	bg.bot.SendMessage(msg, messageDTO)
}

func (bg BotGateway) SendMyAppointmentsMessage(ctx context.Context, user entity.User, appointments []appointment.Appointment, messageDTO tg.MessageDTO, offset int) {
	msgText, keyboard := bg.keyboard.ConfigureGetMyAppointmentsMessage(ctx, user, appointments, offset)
	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	bg.bot.SendMessage(msg, messageDTO)
}

func (bg BotGateway) SendDetailAppointmentMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, appointmentEntity appointment.Appointment) {
	msgText, keyboard := bg.keyboard.ConfigureAppointmentDetailMessage(ctx, user, appointmentEntity)
	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	if msgText == "" {
		bg.SendError(ctx, user, messageDTO)
		return
	}
	msg.ReplyMarkup = keyboard
	bg.bot.SendMessage(msg, messageDTO)
}

func (bg BotGateway) SendEmptyAppointments(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	emptyMessageText, err := bg.messageService.GetMessage(ctx, user, "empty appointments")
	if err != nil {
		bg.SendError(ctx, user, messageDTO)
		return
	}
	msg := tgbotapi.NewMessage(user.GetTgId(), emptyMessageText)
	bg.bot.SendMessage(msg, messageDTO)
}

func (bg BotGateway) SendEmptySchedulesHomeVisit(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	emptyMessageText, err := bg.messageService.GetMessage(ctx, user, "empty schedules home")
	if err != nil {
		bg.SendError(ctx, user, messageDTO)
		return
	}
	msg := tgbotapi.NewMessage(user.GetTgId(), emptyMessageText)
	bg.bot.SendMessage(msg, messageDTO)
}

func (bg BotGateway) SendEmptySchedulePeriods(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	emptyMessageText, err := bg.messageService.GetMessage(ctx, user, "empty schedule periods")
	if err != nil {
		bg.SendError(ctx, user, messageDTO)
		return
	}
	msg := tgbotapi.NewMessage(user.GetTgId(), emptyMessageText)
	bg.bot.SendMessage(msg, messageDTO)
}

func (bg BotGateway) SendWaitMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, waitMessage string) int {
	msgText, _ := bg.messageService.GetMessage(ctx, user, waitMessage)
	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	messageId := bg.bot.SendMessageAndGetId(msg, messageDTO)
	return messageId
}

func (bg BotGateway) SendCalendarMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO, year int, month time.Month, schedulesMap map[time.Time]bool) {
	msgText, keyboard := bg.keyboard.GenerateCalendarKeyboard(ctx, user, year, month, schedulesMap)
	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	bg.bot.SendMessage(msg, messageDTO)
}

func (bg BotGateway) SendChooseAppointmentMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	msgText, keyboard := bg.keyboard.ConfigureChooseAppointmentMessage(ctx, user)
	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	bg.bot.SendMessage(msg, messageDTO)
}

func (bg BotGateway) SendDoctorOrReasonMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	msgText, keyboard := bg.keyboard.ConfigureDoctorOrReasonMessage(ctx, user)
	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	bg.bot.SendMessage(msg, messageDTO)
}

func (bg BotGateway) SendHomeDoctorSpecificationMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO) {
	msgText, keyboard := bg.keyboard.ConfigureChooseHomeDoctorSpecificationMessage(ctx, user)
	msg := tgbotapi.NewMessage(user.GetTgId(), msgText)
	msg.ReplyMarkup = keyboard
	bg.bot.SendMessage(msg, messageDTO)
}

package fast_appointment

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
)

type appointmentService interface {
	GetFastAppointmentSchedules(ctx context.Context) (randomDoctors map[int]appointment.Schedule)
	CreateDraftAppointment(ctx context.Context, tgId int64)
}

type userService interface {
	GetUser(ctx context.Context, tgId int64) (user entity.User, err error)
}

type messageService interface {
	GetMessage(ctx context.Context, user entity.User, name string) (messageText string, err error)
}

type botService interface {
	ConfigureFastAppointmentMessage(
		ctx context.Context,
		userEntity entity.User,
		schedulesMap map[int]appointment.Schedule,
	) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)
}

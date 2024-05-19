package create_appointment

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
)

type appointmentService interface {
	// speciality interfaces in service and gateway
	GetSpecialities(ctx context.Context) (specialities []appointment.Speciality, err error)
	GetTranslatedSpecialities(ctx context.Context, user entity.User, specialities []appointment.Speciality, offset int) (translatedSpecialities map[int]string, unTranslatedSpecialities []string, err error)
	CreateDraftAppointment(ctx context.Context, tgId int64)
}

type userService interface {
	GetUser(ctx context.Context, tgId int64) (user entity.User, err error)
}

type messageService interface {
	GetMessage(ctx context.Context, user entity.User, name string) (messageText string, err error)
}

type botService interface {
	ConfigureGetSpecialityMessage(
		ctx context.Context,
		userEntity entity.User,
		translatedSpecialities map[int]string,
		offset int,
	) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)
}

package text_message

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
)

type appointmentService interface {
	// appointmeent interfaces in service and gateway
	CreateAppointment(ctx context.Context, user entity.User, callbackData string) (appointmentId *int)

	// user interfaces in service and gateway
	GetPatient(ctx context.Context, user entity.User) (result bool)
	CreatePatient(ctx context.Context, user entity.User) (result bool)
}

type userService interface {
	GetUser(ctx context.Context, tgId int64) (user entity.User, err error)
	UpdatePatientId(ctx context.Context, user entity.User, patientId *int) (err error)
	UpdateBirthDate(ctx context.Context, dto tg.TgUserDTO, birthDate string) (user entity.User, err error)
	UpdateThirdName(ctx context.Context, dto tg.TgUserDTO, thirdName string) (user entity.User, err error)
	UpdatePhone(ctx context.Context, dto tg.TgUserDTO, phone string) (user entity.User, err error)
}

type messageService interface {
	GetMessage(ctx context.Context, user entity.User, name string) (messageText string, err error)
}

type botService interface {
	ConfigureChangeLanguageMessage(ctx context.Context, user entity.User) (msgText string, keyboard tgbotapi.InlineKeyboardMarkup)
}
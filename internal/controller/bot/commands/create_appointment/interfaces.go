package create_appointment

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
)

type AppointmentService interface {
	// speciality interfaces in service and gateway
	GetSpecialities(ctx context.Context) (specialities []appointment.Speciality, err error)
	GetTranslatedSpecialities(ctx context.Context, user entity.User, specialities []appointment.Speciality) (translatedSpecialities map[int]string, unTranslatedSpecialities []string, err error)
}

type UserService interface {
	GetUser(ctx context.Context, dto tg.TgUserDTO) (user entity.User, err error)
}

type MessageService interface {
	GetMessage(ctx context.Context, user entity.User, name string) (messageText string, err error)
	SaveMessageLog(ctx context.Context, message tg.MessageDTO) (err error)
}

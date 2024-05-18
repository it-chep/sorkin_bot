package tasks

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
)

type AppointmentService interface {
	GetSpecialities(ctx context.Context) (specialities []appointment.Speciality, err error)
	GetTranslatedSpecialities(
		ctx context.Context,
		user entity.User,
		specialities []appointment.Speciality,
		offset int,
	) (translatedSpecialities map[int]string, unTranslatedSpecialities []string, err error)
}

type UserService interface {
	GetUser(ctx context.Context, dto tg.TgUserDTO) (user entity.User, err error)
}

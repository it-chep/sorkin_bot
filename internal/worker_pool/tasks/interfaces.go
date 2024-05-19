package tasks

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
	tgEntity "sorkin_bot/internal/domain/entity/tg"
	entity "sorkin_bot/internal/domain/entity/user"
)

type appointmentService interface {
	GetSpecialities(ctx context.Context) (specialities []appointment.Speciality, err error)
	GetTranslatedSpecialities(
		ctx context.Context,
		user entity.User,
		specialities []appointment.Speciality,
		offset int,
	) (translatedSpecialities map[int]string, unTranslatedSpecialities []string, err error)
}

type userService interface {
	GetUser(ctx context.Context, tgId int64) (user entity.User, err error)
}

type messageService interface {
	GetSupportLogs(ctx context.Context, minutes int) (logs []tgEntity.MessageLog, err error)
}

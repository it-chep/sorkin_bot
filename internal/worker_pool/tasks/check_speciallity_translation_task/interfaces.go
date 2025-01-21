package check_speciallity_translation_task

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
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

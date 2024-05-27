package keyboards

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
	userEntity "sorkin_bot/internal/domain/entity/user"
)

type messageService interface {
	GetMessage(ctx context.Context, user userEntity.User, name string) (messageText string, err error)
	GetTranslationsBySlugKeyProfession(ctx context.Context, slug string) (translations map[string]appointment.TranslationEntity, err error)
}

type appointmentService interface {
	GetTranslationString(langCode string, translationEntity appointment.TranslationEntity) (translatedSpeciality string)
	GetDoctorInfo(ctx context.Context, user userEntity.User, doctorId int) (doctorEntity appointment.Doctor, err error)
}

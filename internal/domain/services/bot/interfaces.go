package bot

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
	userEntity "sorkin_bot/internal/domain/entity/user"
)

type messageService interface {
	GetMessage(ctx context.Context, user userEntity.User, name string) (messageText string, err error)
}

type readTranslationStorage interface {
	GetTranslationsBySlugKeyProfession(ctx context.Context, slug string) (translations map[string]appointment.TranslationEntity, err error)
}

type appointmentService interface {
	GetSpecialityTranslate(langCode string, translationEntity appointment.TranslationEntity) (translatedSpeciality string)
}

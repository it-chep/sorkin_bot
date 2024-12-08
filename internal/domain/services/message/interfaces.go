package message

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/tg"
)

type ReadRepo interface {
	GetMessageByName(ctx context.Context, name string) (err error, messageEntity entity.Message)
	GetWeekdaysName(ctx context.Context) (err error, messageEntity []entity.Message)
}

type readLogsRepo interface {
	GetSupportLogsByMinutes(ctx context.Context, minutes int) (logs []entity.MessageLog, err error)
}

type SaveMessageUseCase interface {
	Execute(ctx context.Context, messageLog entity.MessageLog) (err error)
}

type readTranslationStorage interface {
	GetTranslationsBySlugKeyProfession(ctx context.Context, slug string) (translations map[string]appointment.TranslationEntity, err error)
}

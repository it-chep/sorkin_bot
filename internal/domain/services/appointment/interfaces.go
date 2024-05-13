package appointment

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
)

type ReadRepo interface {
	GetTranslationsBySlug(ctx context.Context, slug string) (translations map[string]appointment.TranslationEntity, err error)
}

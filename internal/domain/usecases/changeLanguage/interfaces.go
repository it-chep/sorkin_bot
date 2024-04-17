package changeLanguage

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/user"
)

type WriteRepo interface {
	UpdateUserLanguageCode(ctx context.Context, user entity.User, languageCode string) (err error)
}

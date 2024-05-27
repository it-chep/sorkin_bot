package change_language

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/user"
)

type WriteRepo interface {
	UpdateUserState(ctx context.Context, user entity.User) (err error)
	UpdateUserVarcharField(ctx context.Context, user entity.User, field, value string) (err error)
}

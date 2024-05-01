package update_user_birth_date

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/user"
)

type WriteRepo interface {
	UpdateUserVarcharField(ctx context.Context, user entity.User, field, value string) (err error)
}

package update_user_third_name

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/user"
)

type WriteRepo interface {
	UpdateUserVarcharField(ctx context.Context, user entity.User, field, value string) (err error)
}

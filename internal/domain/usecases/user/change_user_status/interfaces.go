package change_user_status

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/user"
)

type writeRepo interface {
	UpdateUserState(ctx context.Context, user entity.User) (err error)
}

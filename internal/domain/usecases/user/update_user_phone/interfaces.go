package update_user_phone

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/user"
)

type WriteRepo interface {
	UpdateUserPhone(ctx context.Context, user entity.User, phone string) (err error)
}

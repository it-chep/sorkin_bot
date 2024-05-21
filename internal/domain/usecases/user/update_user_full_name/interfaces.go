package update_user_full_name

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/user"
)

type WriteRepo interface {
	UpdateUserFullName(ctx context.Context, user entity.User, name, surname, thirdName string) (err error)
}

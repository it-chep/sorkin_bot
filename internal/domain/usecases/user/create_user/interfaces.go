package create_user

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/user"
)

type WriteRepo interface {
	CreateUser(ctx context.Context, user entity.User) (userId int64, err error)
}

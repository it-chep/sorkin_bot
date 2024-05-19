package state_machine

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/user"
)

type userService interface {
	ChangeState(ctx context.Context, tgId int64, state string) (user entity.User, err error)
}

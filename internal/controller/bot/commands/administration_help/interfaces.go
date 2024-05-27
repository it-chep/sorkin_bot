package administration_help

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/user"
)

type messageService interface {
	GetMessage(ctx context.Context, user entity.User, name string) (messageText string, err error)
}

type userService interface {
	GetUser(ctx context.Context, tgId int64) (user entity.User, err error)
}

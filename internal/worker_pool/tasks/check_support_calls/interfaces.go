package check_support_calls

import (
	"context"
	tgEntity "sorkin_bot/internal/domain/entity/tg"
	entity "sorkin_bot/internal/domain/entity/user"
)

type userService interface {
	GetUser(ctx context.Context, tgId int64) (user entity.User, err error)
}

type messageService interface {
	GetSupportLogs(ctx context.Context, minutes int) (logs []tgEntity.MessageLog, err error)
}

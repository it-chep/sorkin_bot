package message

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/tg"
)

type ReadRepo interface {
	GetMessageByName(ctx context.Context, name string) (err error, messageEntity entity.Message)
}

type readLogsRepo interface {
	GetSupportLogsByMinutes(ctx context.Context, minutes int) (logs []entity.MessageLog, err error)
}

type SaveMessageUseCase interface {
	Execute(ctx context.Context, messageLog entity.MessageLog) (err error)
}

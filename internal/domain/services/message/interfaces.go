package message

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/tg"
)

type ReadRepo interface {
	GetMessageByName(ctx context.Context, name string) (err error, messageEntity entity.Message)
}

type SaveMessageUseCase interface {
	Execute(ctx context.Context, messageLog entity.MessageLog) (err error)
}

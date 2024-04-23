package message

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/tg"
)

type SaveMessageUseCase interface {
	Execute(ctx context.Context, messageLog entity.MessageLog) (err error)
}

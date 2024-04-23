package save_message_log

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/tg"
)

type WriteRepo interface {
	CreateMessageLog(ctx context.Context, messageLog entity.MessageLog) (err error)
}

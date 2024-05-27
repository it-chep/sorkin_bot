package middleware

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
)

type MessageService interface {
	SaveMessageLog(ctx context.Context, dto tg.MessageDTO) (err error)
}

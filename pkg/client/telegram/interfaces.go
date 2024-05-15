package telegram

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
)

type MessageService interface {
	SaveMessageLog(ctx context.Context, messageDTO tg.MessageDTO) (err error)
}

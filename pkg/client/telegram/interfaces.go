package telegram

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
)

type messageService interface {
	SaveMessageLog(ctx context.Context, messageDTO tg.MessageDTO) (err error)
}

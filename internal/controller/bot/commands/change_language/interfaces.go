package change_language

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
)

type userService interface {
	GetUser(ctx context.Context, tgId int64) (user entity.User, err error)
}

type botGateway interface {
	SendChangeLanguageMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
}

package start

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
)

type userService interface {
	RegisterNewUser(ctx context.Context, dto tg.TgUserDTO) (user entity.User, err error)
}

type messageService interface {
	GetMessage(ctx context.Context, user entity.User, name string) (messageText string, err error)
}

type botGateway interface {
	SendStartMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
	SendChangeLanguageMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
}

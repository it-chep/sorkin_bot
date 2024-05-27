package exit

import (
	"context"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/user"
)

type botGateway interface {
	SendStartMessage(ctx context.Context, user entity.User, messageDTO tg.MessageDTO)
}

type userService interface {
	GetUser(ctx context.Context, tgId int64) (user entity.User, err error)
	ChangeState(ctx context.Context, tgId int64, state string) (user entity.User, err error)
}

type draftAppointment interface {
	CleanDraftAppointment(ctx context.Context, tgId int64)
}

package user

import (
	"context"
	entity "sorkin_bot/internal/domain/entity/user"
)

type ReadUserStorage interface {
	GetUserByTgID(ctx context.Context, tgID int64) (user entity.User, err error)
}

type CreateUserUseCase interface {
	Execute(ctx context.Context, user entity.User) (userId int64, err error)
}

type UpdateUserLanguageUseCase interface {
	Execute()
}

type UpdatePhoneLanguageUseCase interface {
	Execute()
}

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

type ChangeLanguageUseCase interface {
	Execute(ctx context.Context, user entity.User, languageCode string) (err error)
}

type ChangeStateUseCase interface {
	Execute(ctx context.Context, user entity.User, state string) (err error)
}

type UpdateUserPhoneUseCase interface {
	Execute(ctx context.Context, user entity.User, phone string) (err error)
}

type UpdateUserPatientIdUseCase interface {
	Execute(ctx context.Context, user entity.User, patientId int64) (err error)
}

type UpdatePhoneLanguageUseCase interface {
	Execute()
}

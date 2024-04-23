package create_user

import (
	"context"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/user"
)

type CreateUserUseCase struct {
	writeRepo WriteRepo
	logger    *slog.Logger
}

func NewCreateUserUseCase(writeRepo WriteRepo, logger *slog.Logger) CreateUserUseCase {
	return CreateUserUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc CreateUserUseCase) Execute(ctx context.Context, user entity.User) (userId int64, err error) {
	return uc.writeRepo.CreateUser(ctx, user)
}

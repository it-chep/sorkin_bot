package change_user_status

import (
	"context"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/user"
)

type ChangeStatusUseCase struct {
	writeRepo writeRepo
	logger    *slog.Logger
}

func NewChangeStatusUseCase(writeRepo writeRepo, logger *slog.Logger) ChangeStatusUseCase {
	return ChangeStatusUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc ChangeStatusUseCase) Execute(ctx context.Context, user entity.User) (err error) {
	return uc.writeRepo.UpdateUserState(ctx, user)
}

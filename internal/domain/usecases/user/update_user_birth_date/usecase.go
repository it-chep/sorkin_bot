package update_user_birth_date

import (
	"context"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/user"
)

type UpdateUserBirthDateUseCase struct {
	writeRepo WriteRepo
	logger    *slog.Logger
}

func NewUpdateUserBirthDateUseCase(writeRepo WriteRepo, logger *slog.Logger) UpdateUserBirthDateUseCase {
	return UpdateUserBirthDateUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc UpdateUserBirthDateUseCase) Execute(ctx context.Context, user entity.User, birthDate string) (err error) {
	return uc.writeRepo.UpdateUserVarcharField(ctx, user, "birth_date", birthDate)
}

package update_user_third_name

import (
	"context"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/user"
)

type UpdateUserThirdNameUseCase struct {
	writeRepo WriteRepo
	logger    *slog.Logger
}

func NewUpdateUserThirdNameUseCase(writeRepo WriteRepo, logger *slog.Logger) UpdateUserThirdNameUseCase {
	return UpdateUserThirdNameUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc UpdateUserThirdNameUseCase) Execute(ctx context.Context, user entity.User, thirdName string) (err error) {
	return uc.writeRepo.UpdateUserVarcharField(ctx, user, "third_name", thirdName)
}

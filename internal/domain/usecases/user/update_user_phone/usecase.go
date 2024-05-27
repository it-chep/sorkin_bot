package update_user_phone

import (
	"context"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/user"
)

type UpdateUserPhoneUseCase struct {
	writeRepo WriteRepo
	logger    *slog.Logger
}

func NewUpdateUserPhoneUseCase(writeRepo WriteRepo, logger *slog.Logger) UpdateUserPhoneUseCase {
	return UpdateUserPhoneUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc UpdateUserPhoneUseCase) Execute(ctx context.Context, user entity.User, phone string) (err error) {
	return uc.writeRepo.UpdateUserVarcharField(ctx, user, "phone", phone)
}

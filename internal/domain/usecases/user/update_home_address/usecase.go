package update_home_address

import (
	"context"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/user"
)

type UpdateUserHomeAddressUseCase struct {
	writeRepo WriteRepo
	logger    *slog.Logger
}

func NewUpdateUserHomeAddressUseCase(writeRepo WriteRepo, logger *slog.Logger) UpdateUserHomeAddressUseCase {
	return UpdateUserHomeAddressUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc UpdateUserHomeAddressUseCase) Execute(ctx context.Context, user entity.User, homeAddress string) (err error) {
	return uc.writeRepo.UpdateUserVarcharField(ctx, user, "home_address", homeAddress)
}

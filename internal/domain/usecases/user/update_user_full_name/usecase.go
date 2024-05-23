package update_user_full_name

import (
	"context"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/user"
	"strings"
)

type UpdateUpdateFullNameUseCase struct {
	writeRepo WriteRepo
	logger    *slog.Logger
}

func NewUpdateFullNameUseCase(writeRepo WriteRepo, logger *slog.Logger) UpdateUpdateFullNameUseCase {
	return UpdateUpdateFullNameUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc UpdateUpdateFullNameUseCase) Execute(ctx context.Context, user entity.User, fullName string) (err error) {
	splitName := strings.Split(fullName, " ")
	name, surname, thirdName := splitName[0], splitName[1], "..."
	return uc.writeRepo.UpdateUserFullName(ctx, user, name, surname, thirdName)
}

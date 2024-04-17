package changeLanguage

import (
	"context"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/user"
)

type ChangeLanguageUseCase struct {
	writeRepo WriteRepo
	logger    *slog.Logger
}

func NewChangeLanguageUseCase(writeRepo WriteRepo, logger *slog.Logger) ChangeLanguageUseCase {
	return ChangeLanguageUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc ChangeLanguageUseCase) Execute(ctx context.Context, user entity.User, languageCode string) (err error) {
	return uc.writeRepo.UpdateUserLanguageCode(ctx, user, languageCode)
}

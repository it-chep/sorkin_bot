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
	// TODO менеджер транзакций от авито
	err = uc.writeRepo.UpdateUserLanguageCode(ctx, user, languageCode)
	if err != nil {
		return err
	}
	err = uc.writeRepo.UpdateUserState(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

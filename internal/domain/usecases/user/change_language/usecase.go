package change_language

import (
	"context"
	"fmt"
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
	op := "sorkin_bot.internal.domain.usecases.bot.change_language.usecase.Execute"
	err = uc.writeRepo.UpdateUserVarcharField(ctx, user, "language_code", languageCode)
	if err != nil {
		uc.logger.Error(fmt.Sprintf("error: %s, place: %s", err, op))
		return err
	}
	err = uc.writeRepo.UpdateUserState(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

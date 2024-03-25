package administrator_help

import (
	"log/slog"
)

type AdministratorHelpUseCase struct {
	logger *slog.Logger
}

func NewAdministratorHelpUseCase(repo string, logger *slog.Logger) AdministratorHelpUseCase {
	return AdministratorHelpUseCase{
		logger: logger,
	}
}

func (u AdministratorHelpUseCase) Execute() {
	return
}

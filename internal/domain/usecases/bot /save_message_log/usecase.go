package save_message_log

import (
	"context"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/tg"
)

type SaveMessageLogUseCase struct {
	writeRepo WriteRepo
	logger    *slog.Logger
}

func NewSaveMessageLogUseCase(writeRepo WriteRepo, logger *slog.Logger) SaveMessageLogUseCase {
	return SaveMessageLogUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc SaveMessageLogUseCase) Execute(ctx context.Context, messageLog entity.MessageLog) (err error) {
	return uc.writeRepo.CreateMessageLog(ctx, messageLog)
}

package message

import (
	"context"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	entity "sorkin_bot/internal/domain/entity/tg"
)

type MessageService struct {
	saveMessageUseCase SaveMessageUseCase
	logger             *slog.Logger
}

func NewMessageService(saveMessageUseCase SaveMessageUseCase, logger *slog.Logger) MessageService {
	return MessageService{
		saveMessageUseCase: saveMessageUseCase,
		logger:             logger,
	}
}

func (ms MessageService) SaveMessageLog(ctx context.Context, message tg.MessageDTO) (err error) {
	// todo add photo and video saving
	messageLog := entity.NewMessageLog(
		message.MessageID,
		message.Chat.ID,
		message.Text,
	)
	return ms.saveMessageUseCase.Execute(ctx, messageLog)
}

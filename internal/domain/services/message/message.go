package message

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	tgEntity "sorkin_bot/internal/domain/entity/tg"
	userEntity "sorkin_bot/internal/domain/entity/user"
)

const ServerError = "500 INTERNAL SERVER ERROR, please call /tech_support"

type MessageService struct {
	saveMessageUseCase SaveMessageUseCase
	readRepo           ReadRepo
	readLogsRepo       readLogsRepo
	logger             *slog.Logger
}

func NewMessageService(saveMessageUseCase SaveMessageUseCase, readRepo ReadRepo, readLogsRepo readLogsRepo, logger *slog.Logger) MessageService {
	return MessageService{
		saveMessageUseCase: saveMessageUseCase,
		readRepo:           readRepo,
		readLogsRepo:       readLogsRepo,
		logger:             logger,
	}
}

func (ms MessageService) GetMessage(ctx context.Context, user userEntity.User, name string) (messageText string, err error) {
	op := "sorkin_bot.internal.domain.services.message.message.GetMessage"
	err, message := ms.readRepo.GetMessageByName(ctx, name)
	if err != nil {
		ms.logger.Error(fmt.Sprintf("400 Message Not Found err: %s, place: %s, message_name: %s", err, op, name))
		return ServerError, err
	}
	translatedMessage, err := ms.translateMessage(user, message)
	if err != nil {
		ms.logger.Error(fmt.Sprintf("400 Message Not Found err: %s, place: %s, message_name: %s", err, op, name))
		return ServerError, err
	}
	return translatedMessage, nil
}

func (ms MessageService) translateMessage(user userEntity.User, message tgEntity.Message) (translatedMessage string, err error) {
	languageCode := user.GetLanguageCode()
	switch languageCode {
	case "RU":
		return message.GetRuText(), nil
	case "EN":
		return message.GetEngText(), nil
	case "PT":
		return message.GetRtBRText(), nil
	}
	return ServerError, errors.New(ServerError)
}

func (ms MessageService) SaveMessageLog(ctx context.Context, message tg.MessageDTO) (err error) {
	// todo add photo and video saving
	messageLog := tgEntity.NewMessageLog(
		message.MessageID,
		message.Chat.ID,
		message.Text,
	)
	return ms.saveMessageUseCase.Execute(ctx, messageLog)
}

func (ms MessageService) GetSupportLogs(ctx context.Context, minutes int) (logs []tgEntity.MessageLog, err error) {
	logs, err = ms.readLogsRepo.GetSupportLogsByMinutes(ctx, minutes)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

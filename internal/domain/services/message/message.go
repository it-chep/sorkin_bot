package message

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sorkin_bot/internal/controller/dto/tg"
	"sorkin_bot/internal/domain/entity/appointment"
	tgEntity "sorkin_bot/internal/domain/entity/tg"
	userEntity "sorkin_bot/internal/domain/entity/user"
)

const ServerError = "500 INTERNAL SERVER ERROR, please call /tech_support"

type MessageService struct {
	saveMessageUseCase     SaveMessageUseCase
	readRepo               ReadRepo
	readLogsRepo           readLogsRepo
	readTranslationStorage readTranslationStorage
	logger                 *slog.Logger
}

func NewMessageService(saveMessageUseCase SaveMessageUseCase, readRepo ReadRepo, readLogsRepo readLogsRepo, logger *slog.Logger, readTranslationStorage readTranslationStorage) MessageService {
	return MessageService{
		saveMessageUseCase:     saveMessageUseCase,
		readRepo:               readRepo,
		readLogsRepo:           readLogsRepo,
		logger:                 logger,
		readTranslationStorage: readTranslationStorage,
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
	if user.GetLanguageCode() == nil {
		return message.GetEngText(), nil
	}

	languageCode := *user.GetLanguageCode()
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

func (ms MessageService) GetWeekdaysName(ctx context.Context, user userEntity.User) (translatedMessages []string, err error) {
	op := "sorkin_bot.internal.domain.services.message.message.GetMessage"
	err, messages := ms.readRepo.GetWeekdaysName(ctx)
	if err != nil {
		ms.logger.Error(fmt.Sprintf("400 Message Not Found err: %s, place: %s", err, op))
		return []string{ServerError}, err
	}
	for _, weekday := range messages {
		translatedMessage, err := ms.translateMessage(user, weekday)
		if err != nil {
			ms.logger.Error(fmt.Sprintf("400 Message Not Found err: %s, place: %s", err, op))
			return []string{ServerError}, err
		}
		translatedMessages = append(translatedMessages, translatedMessage)
	}

	return translatedMessages, nil
}

func (ms MessageService) SaveMessageLog(ctx context.Context, message tg.MessageDTO) (err error) {
	// todo add photo and video saving
	//messageLog := tgEntity.NewMessageLog(
	//	message.MessageID,
	//	message.Chat.ID,
	//	message.Text,
	//)
	//return ms.saveMessageUseCase.Execute(ctx, messageLog)
	return nil
	//	todo разобраться
}

func (ms MessageService) GetSupportLogs(ctx context.Context, minutes int) (logs []tgEntity.MessageLog, err error) {
	logs, err = ms.readLogsRepo.GetSupportLogsByMinutes(ctx, minutes)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (ms MessageService) GetTranslationsBySlugKeyProfession(ctx context.Context, slug string) (translations map[string]appointment.TranslationEntity, err error) {
	return ms.readTranslationStorage.GetTranslationsBySlugKeyProfession(ctx, slug)
}

package read_repo

import (
	"log/slog"
	"sorkin_bot/pkg/client/postgres"
)

type TelegramMessageStorage struct {
	client postgres.Client
	logger *slog.Logger
}

func NewTelegramMessageStorage(client postgres.Client, logger *slog.Logger) TelegramMessageStorage {
	return TelegramMessageStorage{
		client: client,
		logger: logger,
	}
}

func (ws TelegramMessageStorage) GetConditionByKeyword() {

}

func (ws TelegramMessageStorage) GetMessageByCondition() {

}

func (ws TelegramMessageStorage) GetButtonByMessageId() {

}

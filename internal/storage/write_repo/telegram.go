package write_repo

import (
	"context"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/tg"
	"sorkin_bot/pkg/client/postgres"
	"time"
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

func (ws TelegramMessageStorage) CreateMessageLog(ctx context.Context, messageLog entity.MessageLog) (err error) {
	op := "internal/storage/write_repo/CreateMessageLog"
	q := `
		insert into message_log (tg_message_id, system_message_id, user_tg_id, time) 
		values ($1, $2, $3, $4);
	`
	ws.logger.Info(op)
	err = ws.client.QueryRow(
		ctx, q, messageLog.GetTgMessageId(), messageLog.GetSystemMessageId(), messageLog.GetUserTgId(), time.Now(),
	).Scan()
	if err != nil {
		return err
	}

	return nil
}

func (ws TelegramMessageStorage) GetConditionByKeyword() {

}

func (ws TelegramMessageStorage) GetMessageByCondition() {

}

func (ws TelegramMessageStorage) GetButtonByMessageId() {

}

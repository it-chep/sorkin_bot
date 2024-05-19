package read_repo

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"log/slog"
	"sorkin_bot/internal/domain/entity/tg"
	"sorkin_bot/internal/storage/dao"
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

func (rs TelegramMessageStorage) GetSupportLogsByMinutes(ctx context.Context, minutes int) (logs []tg.MessageLog, err error) {
	op := "sorkin_bot.internal.storage.read_repo.telegram.GetSupportLogsByMinutes"
	q := fmt.Sprintf(`
		select id, tg_message_id, text, user_tg_id, time
		from message_log
		where text like '/tech_support'
		and time >= now() - interval '%d minutes';
	`, minutes)

	var messageLogDAO []dao.MessageLogDAO
	err = pgxscan.Select(ctx, rs.client, &messageLogDAO, q)
	if err != nil {
		rs.logger.Error(fmt.Sprintf("Error while scanning row: %s op: %s", err, op))
		return logs, err
	}

	for _, m := range messageLogDAO {
		logs = append(logs, m.ToDomain())
	}

	return logs, nil
}

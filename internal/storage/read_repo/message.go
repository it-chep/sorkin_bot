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

type MessageStorage struct {
	client postgres.Client
	logger *slog.Logger
}

func NewReadMessageStorage(client postgres.Client, logger *slog.Logger) MessageStorage {
	return MessageStorage{
		client: client,
		logger: logger,
	}
}

func (mr MessageStorage) GetMessageByName(ctx context.Context, name string) (err error, messageEntity tg.Message) {
	var MessageDAO dao.MessageDAO
	op := "sorkin_bot.internal.storage.read_repo.message.GetMessageByName"
	q := `select id, ru_text, eng_text, pt_br_text from message where name = $1;`

	err = pgxscan.Get(ctx, mr.client, &MessageDAO, q, name)
	if err != nil {
		mr.logger.Error(fmt.Sprintf("error while scanning db rows %s, place: %s, name: %s", err, op, name))
		return err, messageEntity
	}

	messageEntity = MessageDAO.ToDomain()

	return nil, messageEntity
}

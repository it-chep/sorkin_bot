package write_repo

//
//import (
//	"context"
//	"log/slog"
//	entity "sorkin_bot/internal/domain/entity/tg"
//	"sorkin_bot/pkg/client/postgres"
//	"time"
//)
//
//type MessageStorage struct {
//	client postgres.Client
//	logger *slog.Logger
//}
//
//func NewMessageStorage(logger *slog.Logger) MessageStorage {
//	return MessageStorage{
//		logger: logger,
//	}
//}
//
//
//func (ws MessageStorage) CreateMessageLog(ctx context.Context, messageLog entity.MessageLog) (err error) {
//	op := "internal/storage/write_repo/CreateMessageLog"
//	q := `
//		insert into message_log (tg_message_id, system_message_id, user_tg_id, time)
//		values ($1, $2, $3, $4);
//	`
//	ws.logger.Info(op)
//	err = ws.client.QueryRow(
//		ctx, q, messageLog.GetTgMessageId(), messageLog.GetSystemMessageId(), messageLog.GetUserTgId(), time.Now(),
//	).Scan()
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

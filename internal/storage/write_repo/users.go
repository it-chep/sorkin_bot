package write_repo

import (
	"context"
	"fmt"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/pkg/client/postgres"
	"time"
)

type UserStorage struct {
	client postgres.Client
	logger *slog.Logger
}

func NewUserStorage(client postgres.Client, logger *slog.Logger) UserStorage {
	return UserStorage{
		client: client,
		logger: logger,
	}
}

func (ws UserStorage) CreateUser(ctx context.Context, user entity.User) (userID int64, err error) {
	op := "internal/storage/write_repo/CreateUser"
	q := `
		insert into tg_users (tg_id, name, surname, username, registration_time) 
		values ($1, $2, $3, $4, $5, $6, $7, $8) returning id;
	`
	ws.logger.Info(op)
	err = ws.client.QueryRow(
		ctx, q, user.GetTgId(), user.GetFirstName(), user.GetLastName(), user.GetUsername(), time.Now(),
	).Scan(&userID)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func (ws UserStorage) UpdateUserLanguageCode(ctx context.Context, user entity.User, languageCode string) (err error) {
	op := "internal/storage/write_repo/UpdateUserLanguageCode"
	q := `
		update tg_users set language_code = $1 where tg_id = $2;
	`
	ws.logger.Info(op)
	_, err = ws.client.Exec(ctx, q, languageCode, user.GetTgId())
	if err != nil {
		ws.logger.Error(fmt.Sprintf("%s op %s", err, op))
		return err
	}
	return nil
}

func (ws UserStorage) UpdateUserState(ctx context.Context, user entity.User) (err error) {
	op := "internal/storage/write_repo/UpdateUserState"
	q := `
		update tg_users set last_state = $1 where tg_id = $2;
	`

	ws.logger.Info(op)

	_, err = ws.client.Exec(ctx, q, user.GetState(), user.GetTgId())
	if err != nil {
		ws.logger.Error(fmt.Sprintf("%s op %s", err, op))
		return err
	}
	return nil
}

func (ws UserStorage) UpdateUserPhone(ctx context.Context, user entity.User, phone string) (err error) {
	op := "internal/storage/write_repo/UpdateUserPhone"
	q := `
		update tg_users set phone = $1 where tg_id = $2;
	`
	ws.logger.Info(op)
	_, err = ws.client.Exec(ctx, q, phone, user.GetTgId())
	if err != nil {
		ws.logger.Error(fmt.Sprintf("%s op %s", err, op))
		return err
	}
	return nil
}

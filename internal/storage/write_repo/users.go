package write_repo

import (
	"context"
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
		insert into tg_users (tg_id, name, surname, username, registration_time, last_state, phone, language_code) 
		values ($1, $2, $3, $4, $5, $6, $7, $8) returning id;
	`
	ws.logger.Info(op)
	err = ws.client.QueryRow(
		ctx, q, user.GetTgId(), user.GetFirstName(), user.GetLastName(), user.GetUsername(), time.Now(), nil, nil, nil,
	).Scan(&userID)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func (ws UserStorage) UpdateUserLanguageCode(ctx context.Context, user entity.User, newLanguage string) (err error) {
	op := "internal/storage/write_repo/UpdateUserState"
	q := `
		update tg_users set last_state = $1 where tg_id = $2;
	`
	ws.logger.Info(op)
	err = ws.client.QueryRow(ctx, q, newLanguage, user.GetTgId()).Scan()
	if err != nil {
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

	err = ws.client.QueryRow(ctx, q, user.GetState(), user.GetTgId()).Scan()
	if err != nil {
		return err
	}
	return nil
}

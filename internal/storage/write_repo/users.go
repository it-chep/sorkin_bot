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
		insert into tg_users (tg_id, name, surname, username, registration_time, last_state) 
		values ($1, $2, $3, $4, $5, '') returning id;
	`
	// Получаем текущее время
	currentTimeUTC := time.Now().UTC()

	// Получаем время в лисабоне(возможно хардкод, но гибкости пока не требуется)
	location, err := time.LoadLocation("Europe/Lisbon")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	// Преобразуем текущее время в локацию GMT+1
	currentTime := currentTimeUTC.In(location)
	registrationTime := fmt.Sprintf("%02d.%02d.%d %02d:%02d", currentTime.Day(), currentTime.Month(), currentTime.Year(), currentTime.Hour(), currentTime.Minute())

	err = ws.client.QueryRow(
		ctx, q, user.GetTgId(), user.GetFirstName(), user.GetLastName(), user.GetUsername(), registrationTime,
	).Scan(&userID)
	if err != nil {
		ws.logger.Error(fmt.Sprintf("%s op %s", err, op))
		return -1, err
	}

	return userID, nil
}

func (ws UserStorage) UpdateUserPatientId(ctx context.Context, user entity.User, patientId int) (err error) {
	op := "internal/storage/write_repo/UpdateUserPatientId"
	q := `
		update tg_users set patient_id = $1 where tg_id = $2;
	`

	_, err = ws.client.Exec(ctx, q, patientId, user.GetTgId())
	if err != nil {
		ws.logger.Error(fmt.Sprintf("%s op %s", err, op))
		return err
	}

	return nil
}

func (ws UserStorage) UpdateUserState(ctx context.Context, user entity.User) (err error) {
	op := "internal.storage.write_repo.UpdateUserState"
	q := `
		update tg_users set last_state = $1 where tg_id = $2;
	`

	_, err = ws.client.Exec(ctx, q, user.GetState(), user.GetTgId())
	if err != nil {
		ws.logger.Error(fmt.Sprintf("%s op %s", err, op))
		return err
	}
	return nil
}

func (ws UserStorage) UpdateUserVarcharField(ctx context.Context, user entity.User, field, value string) (err error) {
	op := "internal.storage.write_repo.UpdateUserVarcharField"
	q := fmt.Sprintf(`
		update tg_users set %s = $1 where tg_id = $2;
	`, field)

	_, err = ws.client.Exec(ctx, q, value, user.GetTgId())
	if err != nil {
		ws.logger.Error(fmt.Sprintf("%s op %s", err, op))
		return err
	}
	return nil
}

func (ws UserStorage) UpdateUserFullName(ctx context.Context, user entity.User, name, surname, thirdName string) (err error) {
	op := "internal.storage.write_repo.UpdateUserFullName"
	q := `update tg_users set name = $1, surname = $2, third_name = $3 where tg_id = $4;`

	_, err = ws.client.Exec(ctx, q, name, surname, thirdName, user.GetTgId())
	if err != nil {
		ws.logger.Error(fmt.Sprintf("%s op %s", err, op))
		return err
	}
	return nil
}

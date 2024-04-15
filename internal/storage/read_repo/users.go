package read_repo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/storage/dao"
	"sorkin_bot/pkg/client/postgres"
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

func (u UserStorage) GetUserByTgID(ctx context.Context, userID int64) (user entity.User, err error) {
	q := "select tg_id, name, surname, username, last_state, phone, language_code from tg_users where tg_id = $1"

	var userDAO dao.UserDAO

	err = u.client.QueryRow(ctx, q, userID).Scan(
		&userDAO.TgId,
		&userDAO.FirstName,
		&userDAO.LastName,
		&userDAO.Username,
		&userDAO.LastState,
		&userDAO.Phone,
		&userDAO.LanguageCode,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, nil
		}
		return entity.User{}, err
	}

	// Create and return a new User entity
	user = *userDAO.ToDomain()
	return user, nil
}

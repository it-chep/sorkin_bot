package read_repo

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"log/slog"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/internal/storage/dao"
	"sorkin_bot/pkg/client/postgres"
)

// Todo use scany

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

func (rs UserStorage) GetUserByTgID(ctx context.Context, userID int64) (user entity.User, err error) {
	op := "internal/storage/read_repo/users/GetUserByTgID"
	q := "select tg_id, name, surname, username, last_state, phone, language_code from tg_users where tg_id = $1"

	var userDAO dao.UserDAO
	rs.logger.Info(op)
	err = pgxscan.Get(ctx, rs.client, &userDAO, q, userID)
	if err != nil {
		rs.logger.Error(fmt.Sprintf("Error while scanning row: %s", err))
		return entity.User{}, err
	}

	// Create and return a new User entity
	user = *userDAO.ToDomain()
	return user, nil
}

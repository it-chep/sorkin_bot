package dao

import entity "sorkin_bot/internal/domain/entity/user"

type UserDAO struct {
	TgId         int64  `db:"tg_id"`
	FirstName    string `db:"name"`
	LastName     string `db:"surname"`
	Username     string `db:"username"`
	LanguageCode string `db:"language_code"`
	Phone        string `db:"phone"`
	LastState    string `db:"last_state"`
}

func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

func (dao *UserDAO) ToDomain() *entity.User {
	return entity.NewUser(
		dao.TgId,
		dao.FirstName,
	)
}

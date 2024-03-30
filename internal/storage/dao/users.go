package dao

import entity "sorkin_bot/internal/domain/entity/user"

type UserDAO struct {
	TgId         int64  `sql:"tg_id"`
	FirstName    string `sql:"name"`
	LastName     string `sql:"surname"`
	Username     string `sql:"username"`
	LanguageCode string `sql:"language_code"`
	Phone        string `sql:"phone"`
	LastState    string `sql:"last_state"`
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

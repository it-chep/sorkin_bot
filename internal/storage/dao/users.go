package dao

import (
	entity "sorkin_bot/internal/domain/entity/user"
)

type UserDAO struct {
	TgId             int64   `db:"tg_id"`
	FirstName        string  `db:"name"`
	LastName         *string `db:"surname"`
	Username         *string `db:"username"`
	LanguageCode     *string `db:"language_code"`
	Phone            *string `db:"phone"`
	LastState        *string `db:"last_state"`
	PatientId        *int    `db:"patient_id"`
	RegistrationTime string  `db:"registration_time"`
	BirthDate        *string `db:"birth_date"`
	HomeAddress      *string `db:"home_address"`
}

func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

func (dao *UserDAO) ToDomain() *entity.User {
	return entity.NewUser(
		dao.TgId,
		dao.FirstName,
		entity.WithUsrLanguageCode(dao.LanguageCode),
		entity.WithUsrState(dao.LastState),
		entity.WithUsrUsername(dao.Username),
		entity.WithUsrLastName(dao.LastName),
		entity.WithUsrPhone(dao.Phone),
		entity.WithUsrPatientId(dao.PatientId),
		entity.WithRegistrationTime(dao.RegistrationTime),
		entity.WithBirthDate(dao.BirthDate),
		entity.WithUsrHomeAddress(dao.HomeAddress),
	)
}

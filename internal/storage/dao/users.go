package dao

import (
	"database/sql"
	entity "sorkin_bot/internal/domain/entity/user"
)

type UserDAO struct {
	TgId             int64          `db:"tg_id"`
	FirstName        string         `db:"name"`
	LastName         string         `db:"surname"`
	Username         string         `db:"username"`
	LanguageCode     string         `db:"language_code"`
	Phone            string         `db:"phone"`
	LastState        string         `db:"last_state"`
	PatientId        *int           `db:"patient_id"`
	RegistrationTime string         `db:"registration_time"`
	BirthDate        sql.NullString `db:"birth_date"`
}

func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

func (dao *UserDAO) ToDomain() *entity.User {
	patientID := -1
	if dao.PatientId != nil {
		patientID = *dao.PatientId
	}

	return entity.NewUser(
		dao.TgId,
		dao.FirstName,
		entity.WithUsrLanguageCode(dao.LanguageCode),
		entity.WithUsrState(dao.LastState),
		entity.WithUsrUsername(dao.Username),
		entity.WithUsrLastName(dao.LastName),
		entity.WithUsrPhone(dao.Phone),
		entity.WithUsrPatientId(patientID),
		entity.WithRegistrationTime(dao.RegistrationTime),
		entity.WithBirthDate(dao.BirthDate.String),
	)
}

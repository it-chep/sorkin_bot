package dao

import "sorkin_bot/internal/domain/entity/appointment"

type AppointmentDAO struct {
	TimeStart    *string `db:"time_start,omitempty"`
	TimeEnd      *string `db:"time_end,omitempty"`
	DoctorId     *int    `db:"doctor_id,omitempty"`
	TgId         *int64  `db:"tg_id,omitempty"`
	SpecialityId *int    `db:"speciality_id,omitempty"`
	Date         *string `db:"date,omitempty"`
}

func (a *AppointmentDAO) ToDomain() appointment.DraftAppointment {
	return appointment.NewDraftAppointment(
		a.SpecialityId,
		a.DoctorId,
		a.TgId,
		a.TimeStart,
		a.TimeEnd,
		a.Date,
	)
}

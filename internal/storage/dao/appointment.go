package dao

import "sorkin_bot/internal/domain/entity/appointment"

type AppointmentDAO struct {
	id           *int    `db:"id"`
	timeStart    *string `db:"time_start,omitempty"`
	timeEnd      *string `db:"time_end,omitempty"`
	doctorId     *int    `db:"doctor_id,omitempty"`
	tgId         *int64  `db:"tg_id,omitempty"`
	specialityId *int    `db:"speciality_id,omitempty"`
}

func (a *AppointmentDAO) ToDomain() appointment.DraftAppointment {
	return appointment.NewDraftAppointment(
		*a.id,
		a.specialityId,
		a.doctorId,
		a.tgId,
		a.timeStart,
		a.timeEnd,
	)
}

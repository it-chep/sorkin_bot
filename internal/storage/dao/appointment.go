package dao

import "sorkin_bot/internal/domain/entity/appointment"

type AppointmentDAO struct {
	TimeStart       *string `db:"time_start,omitempty"`
	TimeEnd         *string `db:"time_end,omitempty"`
	DoctorId        *int    `db:"doctor_id,omitempty"`
	DoctorName      *string `db:"doctor_name,omitempty"`
	TgId            *int64  `db:"tg_id,omitempty"`
	SpecialityId    *int    `db:"speciality_id,omitempty"`
	Date            *string `db:"date,omitempty"`
	AppointmentType *string `db:"type,omitempty"`
}

func (a *AppointmentDAO) ToDomain() appointment.DraftAppointment {
	var appointmentType *appointment.AppointmentType
	if a.AppointmentType != nil {
		value := appointment.AppointmentType(*a.AppointmentType)
		appointmentType = &value
	}
	return appointment.NewDraftAppointment(
		a.SpecialityId,
		a.DoctorId,
		a.TgId,
		a.DoctorName,
		a.TimeStart,
		a.TimeEnd,
		a.Date,
		appointmentType,
	)
}

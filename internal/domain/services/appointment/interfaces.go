package appointment

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
)

type Appointment interface {
	FastAppointment(ctx context.Context)
	CreateAppointment(ctx context.Context) (err error, appointmentId int)
	CancelAppointment(ctx context.Context, movedTo string, appointmentId int) (err error, result bool)
	ConfirmAppointment(ctx context.Context, appointmentId int) (err error, result bool)
	RescheduleAppointment(ctx context.Context, movedTo string) (err error)
	MyAppointments(ctx context.Context) (err error, appointments []appointment.Appointment)
	DetailAppointment(ctx context.Context) (err error, appointmentEntity appointment.Appointment)
	GetDoctors(ctx context.Context) (err error, doctors []appointment.Doctor)
	GetSpecialities(ctx context.Context) (err error, specialities []appointment.Speciality)
	//GetSchedules(ctx context.Context) (err error, response mis_dto.GetSpecialityResponse)
}

type ReadRepo interface {
	GetTranslationsBySlug(ctx context.Context, slug string) (translations map[string]appointment.TranslationEntity, err error)
}

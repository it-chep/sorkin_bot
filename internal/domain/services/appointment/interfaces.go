package appointment

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
	userEntity "sorkin_bot/internal/domain/entity/user"
)

type Appointment interface {
	FastAppointment(ctx context.Context) (err error, schedulesMap map[int][]appointment.Schedule)
	CreateAppointment(ctx context.Context, userEntity userEntity.User, doctorId int, timeStart, timeEnd string) (err error, appointmentId int)
	CancelAppointment(ctx context.Context, movedTo string, appointmentId int) (err error, result bool)
	ConfirmAppointment(ctx context.Context, appointmentId int) (err error, result bool)
	MyAppointments(ctx context.Context, user userEntity.User) (err error, appointments []appointment.Appointment)
	DetailAppointment(ctx context.Context, user userEntity.User, appointmentId int) (err error, appointmentEntity appointment.Appointment)
	GetDoctors(ctx context.Context, specialityId int) (err error, doctors []appointment.Doctor)
	GetSpecialities(ctx context.Context) (err error, specialities []appointment.Speciality)
	GetPatientById(ctx context.Context, patientId int) (err error)
	CreatePatient(ctx context.Context, user userEntity.User) (err error, patientId int)
	GetSchedules(ctx context.Context, doctorId int, timeStart string) (err error, schedulesMap map[int][]appointment.Schedule)
}

type ReadRepo interface {
	GetTranslationsBySlug(ctx context.Context, slug string) (translations map[string]appointment.TranslationEntity, err error)
}

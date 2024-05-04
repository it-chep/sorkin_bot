package appointment

import (
	"context"
	"sorkin_bot/internal/clients/gateways/dto"
	"sorkin_bot/internal/domain/entity/appointment"
)

type SchedulesActions interface {
	GetSchedules(ctx context.Context, doctorId int, timeStart string) (schedulesMap map[int][]dto.ScheduleDTO, err error)
	FastAppointment(ctx context.Context) (schedulesMap map[int][]dto.ScheduleDTO, err error)
}

type AppointmentsActions interface {
	CreateAppointment(ctx context.Context, patientId, doctorId int, timeStart, timeEnd string) (appointmentId *int, err error)
	CancelAppointment(ctx context.Context, movedTo string, appointmentId int) (result bool, err error)
	ConfirmAppointment(ctx context.Context, appointmentId int) (result bool, err error)
	MyAppointments(ctx context.Context, patientId int, registrationTime string) (appointments []dto.AppointmentDTO, err error)
	DetailAppointment(ctx context.Context, patientId, appointmentId int, registrationTime string) (appointmentDTO dto.AppointmentDTO, err error)
}

type DoctorsActions interface {
	GetDoctors(ctx context.Context, specialityId int) (doctors []dto.DoctorDTO, err error)
}

type PatientActions interface {
	GetPatientById(ctx context.Context, patientId int) (patientDTO dto.CreatedPatientDTO, err error)
	CreatePatient(ctx context.Context, userDTO dto.PatientDTO) (patientId *int, err error)
}

type SpecialityActions interface {
	GetSpecialities(ctx context.Context) (specialities []dto.SpecialityDTO, err error)
}

type ReadRepo interface {
	GetTranslationsBySlug(ctx context.Context, slug string) (translations map[string]appointment.TranslationEntity, err error)
}

package adapter

import (
	"context"
	"sorkin_bot/internal/clients/gateways/dto"
	entity "sorkin_bot/internal/domain/entity/user"
)

type Gateway interface {
	SchedulesActions
	AppointmentsActions
	DoctorsActions
	PatientActions
	SpecialityActions
}

type SchedulesActions interface {
	GetSchedulesByDoctorId(ctx context.Context, doctorId int, timeStart, timeEnd string) (schedulesMap map[int][]dto.ScheduleDTO, err error)
	GetAvailableDoctorIdsFromSchedulePeriods(ctx context.Context, doctorIdsBySpeciality []int, timeStart, timeEnd string) (availableDoctorIds []int, err error)
	GetSchedulesManyDoctors(ctx context.Context, doctorIds []int, timeStart, timeEnd string) (schedules []dto.ScheduleDTO, err error)
	GetSchedulePeriodsByDoctorId(ctx context.Context, doctorId int, timeStart, timeEnd string) (schedulePeriods []dto.SchedulePeriodDTO, err error)
}

type AppointmentsActions interface {
	CreateAppointment(ctx context.Context, createAppointmentDTO dto.CreateAppointmentDTO) (appointmentId *int, err error)
	CancelAppointment(ctx context.Context, movedTo string, appointmentId int) (result bool, err error)
	ConfirmAppointment(ctx context.Context, appointmentId int) (result bool, err error)
	MyAppointments(ctx context.Context, patientId int, registrationTime string) (appointments []dto.AppointmentDTO, err error)
	DetailAppointment(ctx context.Context, patientId, appointmentId int, registrationTime string) (appointmentDTO dto.AppointmentDTO, err error)
}

type DoctorsActions interface {
	GetDoctorsBySpecialityId(ctx context.Context, specialityId int) (doctors []dto.DoctorDTO, err error)
	GetDoctorInfo(ctx context.Context, doctorId int) (doctorDTO dto.DoctorDTO, err error)
	GetDoctors(ctx context.Context, homeVisit, onlineAppointment, clinicAppointment bool) (doctors []dto.DoctorDTO, err error)
}

type PatientActions interface {
	GetPatientById(ctx context.Context, patientId int) (patientDTO dto.CreatedPatientDTO, err error)
	CreatePatient(ctx context.Context, userDTO dto.PatientDTO) (patientId *int, err error)
	GetPatientByBirthDate(ctx context.Context, user entity.User) (patientDTO dto.CreatedPatientDTO, err error)
}

type SpecialityActions interface {
	GetSpecialities(ctx context.Context) (specialities []dto.SpecialityDTO, err error)
}

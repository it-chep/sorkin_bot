package adapter

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
)

type Cache interface {
	//	todo cache
}

// todo add adapter on gateway
// todo use cache

type AppointmentServiceAdapter struct {
	cache   Cache
	gateway Gateway
}

func NewAppointmentServiceAdapter() AppointmentServiceAdapter {
	return AppointmentServiceAdapter{}
}

func (a *AppointmentServiceAdapter) GetSpecialities(ctx context.Context) (specialities []appointment.Speciality) {
	cachedSpecialities, ok := a.cache.Get(ctx, "specialities")
	if !ok {
		specialities, err := a.gateway.GetSpecialities(ctx)
		if err != nil {
			return
		}
	}
	for _, speciality := range cachedSpecialities {
		specialities = append(specialities, speciality.ToEntity())
	}
	return specialities
}

func (a *AppointmentServiceAdapter) GetDoctors(ctx context.Context, specialityId int) (doctors []appointment.Doctor) {
	cachedDoctors, ok := a.cache.Get(ctx, "doctors")
	if !ok {
		doctors, err := a.gateway.GetDoctors(ctx, specialityId)
		if err != nil {
			return
		}
	}
	for _, doctor := range doctors {
		doctors = append(doctors, doctor.ToEntity())
	}
	return doctors
}

func (a *AppointmentServiceAdapter) MyAppointments(ctx context.Context, user entity.User) (appointments []appointment.Appointment) {
	cachedAppointments, ok := a.cache.Get(ctx, "appointments")
	if user.GetPatientId() == nil {
		return
	}
	if !ok {
		appointments, err := a.gateway.MyAppointments(ctx, *user.GetPatientId(), user.GetRegistrationTime())
		if err != nil {
			return
		}
	}
	for _, appointmentDTO := range appointments {
		appointments = append(appointments, appointmentDTO.ToEntity())
	}
	return appointments
}

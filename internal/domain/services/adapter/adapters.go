package adapter

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"sorkin_bot/pkg/client/inmemory_cache"
	"time"
)

// todo add adapter on gateway
// todo use cache

type AppointmentServiceAdapter struct {
	cache   inmemory_cache.Cache[string, any]
	gateway Gateway
}

func NewAppointmentServiceAdapter() AppointmentServiceAdapter {
	return AppointmentServiceAdapter{
		cache: *inmemory_cache.NewCache[string, any](time.Second * 10),
	}
}

func (a *AppointmentServiceAdapter) GetSpecialities(ctx context.Context) (specialities []appointment.Speciality) {
	cachedSpecialities, ok := a.cache.Get("specialities")
	if !ok {
		specialities, err := a.gateway.GetSpecialities(ctx)
		if err != nil {
			return
		}
	}
	for _, speciality := range specialities {
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

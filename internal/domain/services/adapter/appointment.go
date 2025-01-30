package adapter

import (
	"context"
	"fmt"
	"sorkin_bot/internal/clients/gateways/dto"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"time"
)

func (a *AppointmentServiceAdapter) CreateAppointment(ctx context.Context, user entity.User, draftAppointment appointment.DraftAppointment, doctorId int, timeStart, timeEnd string) (appointmentId *int, err error) {
	createAppointmentDTO := dto.CreateAppointmentDTO{
		PatientId:         *user.GetPatientId(),
		DoctorId:          doctorId,
		TimeStart:         timeStart,
		TimeEnd:           timeEnd,
		HomeAddress:       user.GetHomeAddress(),
		HomeVisit:         false,
		OnlineAppointment: false,
		ClinicAppointment: false,
	}

	switch *draftAppointment.GetAppointmentType() {
	case appointment.HomeAppointment:
		createAppointmentDTO.HomeVisit = true
	case appointment.ClinicAppointment:
		createAppointmentDTO.ClinicAppointment = true
	case appointment.OnlineAppointment:
		createAppointmentDTO.OnlineAppointment = true
	}

	appointmentId, err = a.gateway.CreateAppointment(ctx, createAppointmentDTO)
	if err != nil {
		return nil, err
	}

	a.cache.Del(fmt.Sprintf("%d_appointments", *user.GetPatientId()))
	return appointmentId, nil
}

func (a *AppointmentServiceAdapter) MyAppointments(ctx context.Context, user entity.User) (appointments []appointment.Appointment) {
	if user.GetPatientId() == nil {
		return
	}

	cacheKey := fmt.Sprintf("%d_appointments", *user.GetPatientId())
	cachedAppointments, ok := a.cache.Get(cacheKey)

	if !ok || cachedAppointments == nil {
		a.cache.Del(cacheKey)
		appointmentsDTO, err := a.gateway.MyAppointments(ctx, *user.GetPatientId(), user.GetRegistrationTime())
		if err != nil {
			return
		}

		appointments = a.cacheMyAppointments(user, appointmentsDTO)
		return appointments
	}

	appointmentsFromCache, ok := cachedAppointments.([]appointment.Appointment)

	if !ok || len(appointmentsFromCache) == 0 {
		a.cache.Del(cacheKey)
		appointmentsDTO, err := a.gateway.MyAppointments(ctx, *user.GetPatientId(), user.GetRegistrationTime())
		if err != nil {
			return
		}

		appointments = a.cacheMyAppointments(user, appointmentsDTO)
		return appointments
	}

	return appointmentsFromCache
}

func (a *AppointmentServiceAdapter) CancelAppointment(ctx context.Context, user entity.User, appointmentId int) (result bool, err error) {
	a.cache.Del(fmt.Sprintf("%d_appointments", *user.GetPatientId()))
	a.cache.Del(fmt.Sprintf("%d_%d_appointment", *user.GetPatientId(), appointmentId))

	result, err = a.gateway.CancelAppointment(ctx, "", appointmentId)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (a *AppointmentServiceAdapter) ConfirmAppointment(ctx context.Context, appointmentId int) (result bool, err error) {
	result, err = a.gateway.ConfirmAppointment(ctx, appointmentId)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (a *AppointmentServiceAdapter) DetailAppointment(ctx context.Context, user entity.User, appointmentId int) (appointment.Appointment, error) {
	cachedAppointment, ok := a.cache.Get(fmt.Sprintf("%d_%d_appointment", user.GetPatientId(), appointmentId))
	if !ok {
		appointmentDTO, err := a.gateway.DetailAppointment(ctx, *user.GetPatientId(), appointmentId, user.GetRegistrationTime())
		if err != nil {
			return appointment.Appointment{}, err
		}
		a.cache.Set(fmt.Sprintf("%d_%d_appointment", appointmentDTO.PatientId, appointmentDTO.Id), appointmentDTO.ToDomain(), 10*time.Minute)
		return appointmentDTO.ToDomain(), nil
	}
	return cachedAppointment.(appointment.Appointment), nil
}

func (a *AppointmentServiceAdapter) RescheduleAppointment(ctx context.Context, user entity.User, moved_to string, appointmentId int) (err error) {
	return nil
}

func (a *AppointmentServiceAdapter) cacheMyAppointments(user entity.User, appointmentsDTO []dto.AppointmentDTO) (appointments []appointment.Appointment) {
	for _, appointmentDTO := range appointmentsDTO {
		appointments = append(appointments, appointmentDTO.ToDomain())
		key := fmt.Sprintf("%d_%d_appointment", appointmentDTO.PatientId, appointmentDTO.Id)
		a.cache.Set(key, appointmentDTO.ToDomain(), 5*time.Minute)
	}

	a.cache.Set(fmt.Sprintf("%d_appointments", *user.GetPatientId()), appointments, 10*time.Minute)

	return appointments
}

func (a *AppointmentServiceAdapter) GetAppointmentsForNotifying(ctx context.Context, dateFrom, dateTo string) ([]appointment.Appointment, error) {
	_ = "adapter.GetAppointmentsForNotifying"
	var appointments []appointment.Appointment

	appointmentsDTO, err := a.gateway.GetAppointmentsForNotifying(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	for _, appointmentDTO := range appointmentsDTO {
		appointments = append(appointments, appointmentDTO.ToDomain())
	}

	return appointments, nil
}

package adapter

import (
	"context"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"time"
)

func (a *AppointmentServiceAdapter) CreateAppointment(ctx context.Context, user entity.User, doctorId int, timeStart, timeEnd string) (appointmentId *int, err error) {
	//appointmentId, err = a.gateway.CreateAppointment(ctx, *user.GetPatientId(), doctorId, timeStart, timeEnd)
	//if err != nil {
	//	return nil, err
	//}
	fixedAppointmentId := 1

	return &fixedAppointmentId, nil
}

func (a *AppointmentServiceAdapter) MyAppointments(ctx context.Context, user entity.User) (appointments []appointment.Appointment) {
	cachedAppointments, ok := a.cache.Get("appointments")

	if !ok {
		appointmentsDTO, err := a.gateway.MyAppointments(ctx, *user.GetPatientId(), user.GetRegistrationTime())
		if err != nil {
			return
		}

		for _, appointmentDTO := range appointmentsDTO {
			appointments = append(appointments, appointmentDTO.ToDomain())
			key := fmt.Sprintf("%d_%d_appointment", appointmentDTO.PatientId, appointmentDTO.Id)
			a.cache.Set(key, appointmentDTO, 12*time.Hour)
		}
		a.cache.Set("appointments", appointments, 12*time.Hour)
		return appointments
	}

	return cachedAppointments.([]appointment.Appointment)
}

func (a *AppointmentServiceAdapter) CancelAppointment(ctx context.Context, user entity.User, appointmentId int) (result bool, err error) {
	a.cache.Del(fmt.Sprintf("%d_appointments", *user.GetPatientId()))
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
		a.cache.Set(fmt.Sprintf("%d_%d_appointment", appointmentDTO.PatientId, appointmentDTO.Id), appointmentDTO.ToDomain(), 12*time.Hour)
		return appointmentDTO.ToDomain(), nil
	}
	return cachedAppointment.(appointment.Appointment), nil
}

func (a *AppointmentServiceAdapter) RescheduleAppointment(ctx context.Context, user entity.User, moved_to string, appointmentId int) (err error) {
	return nil
}

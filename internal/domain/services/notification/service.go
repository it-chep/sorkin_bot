package notification

import (
	"context"
	"errors"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
)

type Service struct {
	notifyGateway notifyGateway
}

func (s *Service) NotifyCancelAppointment(ctx context.Context, appointment appointment.Appointment) error {
	data, ok := clinicDataMap[appointment.GetClinicId()]
	if !ok {
		return errors.New("invalid clinic id")
	}

	cancelAppointmentMessage := s.prepareMessage(cancelAppointmentTemplate, appointment, data)

	err := s.notifyGateway.SendNotification(ctx, []string{appointment.GetPatientPhone()}, cancelAppointmentMessage)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) NotifySoonAppointment(ctx context.Context, appointment appointment.Appointment) error {
	data, ok := clinicDataMap[appointment.GetClinicId()]
	if !ok {
		return errors.New("invalid clinic id")
	}

	soonAppointmentMessage := s.prepareMessage(visitReminderTemplate, appointment, data)

	err := s.notifyGateway.SendNotification(ctx, []string{appointment.GetPatientPhone()}, soonAppointmentMessage)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) prepareMessage(template string, appointment appointment.Appointment, clinic clinicData) string {
	return fmt.Sprintf(
		template,
		appointment.GetStringDateTimeStart(),
		appointment.GetPatientName(),
		appointment.GetStringDateStart(),
		appointment.GetStringTimeStart(),
		appointment.GetClinic(),
		appointment.GetDoctor(),
		clinic.address,
		clinic.phone,
	)
}

func NewService(notifyGateway notifyGateway) *Service {
	return &Service{
		notifyGateway: notifyGateway,
	}
}

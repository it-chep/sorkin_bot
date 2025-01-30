package notification

import (
	"context"
	"errors"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
	"strings"
)

type Service struct {
	notifyGateway notifyGateway
	misGateway    misGateway
}

func (s *Service) NotifyCancelAppointment(ctx context.Context, appointment appointment.Appointment) error {
	patientPhone, err := s.getPatientPhone(ctx, appointment)
	if err != nil {
		return err
	}

	data, ok := clinicDataMap[appointment.GetClinicId()]
	if !ok {
		return errors.New("invalid clinic id")
	}

	cancelAppointmentMessage := s.prepareMessage(cancelAppointmentTemplate, appointment, data)

	err = s.notifyGateway.SendNotification(ctx, []string{patientPhone}, cancelAppointmentMessage)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) NotifySoonAppointment(ctx context.Context, appointment appointment.Appointment) error {
	patientPhone, err := s.getPatientPhone(ctx, appointment)
	if err != nil {
		return err
	}

	data, ok := clinicDataMap[appointment.GetClinicId()]
	if !ok {
		return errors.New("invalid clinic id")
	}

	soonAppointmentMessage := s.prepareMessage(visitReminderTemplate, appointment, data)

	err = s.notifyGateway.SendNotification(ctx, []string{patientPhone}, soonAppointmentMessage)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) cleanPhoneNumber(phone string) string {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	return phone
}

func (s *Service) getPatientPhone(ctx context.Context, appointment appointment.Appointment) (string, error) {
	patientPhone := appointment.GetPatientPhone()
	if len(patientPhone) == 0 {
		patientDTO, err := s.misGateway.GetPatientById(ctx, appointment.GetPatientId())
		if err != nil {
			return "", err
		}
		patientPhone = patientDTO.Phone
	}

	patientPhone = s.cleanPhoneNumber(patientPhone)

	return patientPhone, nil
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

func NewService(notifyGateway notifyGateway, misGateway misGateway) *Service {
	return &Service{
		notifyGateway: notifyGateway,
		misGateway:    misGateway,
	}
}

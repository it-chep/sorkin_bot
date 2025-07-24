package notification

import (
	"context"
	"errors"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
	"strings"
)

type NotificationType = int

const (
	CreateAppointment NotificationType = iota
	CancelAppointment
	RemindAboutAppointment
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

	// Если визит перенесен, то выходим
	if appointment.MovedToID() != 0 || appointment.MovedFromID() != 0 {
		return nil
	}

	data, ok := clinicDataMap[appointment.ClinicId()]
	if !ok {
		return errors.New("invalid clinic id")
	}

	cancelAppointmentMessage := s.prepareMessage(cancelAppointmentTemplate, appointment, data)
	return s.notifyGateway.SendNotification(ctx, []string{patientPhone}, cancelAppointmentMessage)
}

func (s *Service) NotifyCreateAppointment(ctx context.Context, appointment appointment.Appointment) error {
	patientPhone, err := s.getPatientPhone(ctx, appointment)
	if err != nil {
		return err
	}

	// Если визит перенесен, то выходим
	if appointment.MovedToID() != 0 || appointment.MovedFromID() != 0 {
		return nil
	}

	data, ok := clinicDataMap[appointment.ClinicId()]
	if !ok {
		return errors.New("invalid clinic id")
	}

	message := s.getMessage(CreateAppointment, appointment, data, patientPhone)
	return s.notifyGateway.SendNotification(ctx, []string{patientPhone}, message)
}

func (s *Service) NotifySoonAppointment(ctx context.Context, appointment appointment.Appointment) error {
	patientPhone, err := s.getPatientPhone(ctx, appointment)
	if err != nil {
		return err
	}

	data, ok := clinicDataMap[appointment.ClinicId()]
	if !ok {
		return errors.New("invalid clinic id")
	}

	message := s.getMessage(RemindAboutAppointment, appointment, data, patientPhone)
	return s.notifyGateway.SendNotification(ctx, []string{patientPhone}, message)
}

// теперь getMessage принимает patientPhone
func (s *Service) getMessage(
	notificationType NotificationType,
	appointment appointment.Appointment,
	data clinicData,
	patientPhone string,
) string {
	// Выездной визит
	if appointment.IsOutside() {
		switch notificationType {
		case CreateAppointment:
			return fmt.Sprintf(
				CreateHouseCallAppointmentTemplate,
				appointment.PatientName(),
				appointment.GetStringDateStart(),
				appointment.Doctor(),
				data.phone,
			)
		case RemindAboutAppointment:
			return fmt.Sprintf(
				HomeVisitReminderTemplate,
				appointment.PatientName(),
				appointment.GetStringDateStart(),
				appointment.Doctor(),
				patientPhone,  // номер пациента для ссылки
				data.phone,
			)
		}
	}


	// Телемедицина
	if appointment.IsTelemedicine() {
		switch notificationType {
		case CreateAppointment:
			return fmt.Sprintf(
				CreateOnlineAppointmentTemplate,
				appointment.PatientName(),
				appointment.GetStringDateStart(),
				appointment.GetStringTimeStart(),
				appointment.Doctor(),
				data.phone,
			)
		case RemindAboutAppointment:
			return fmt.Sprintf(
				OnlineVisitReminderTemplate,
				appointment.PatientName(),
				appointment.GetStringDateStart(),
				appointment.GetStringTimeStart(),
				appointment.Doctor(),
				patientPhone,  // номер пациента для ссылки
				data.phone,
			)
		}
	}

	// Приём в клинике
	switch notificationType {
	case CreateAppointment:
		// здесь patientPhone не нужен, поэтому не передаётся
		return fmt.Sprintf(
			CreateInClinicAppointmentTemplate,
			appointment.PatientName(),
			appointment.GetStringDateStart(),
			appointment.GetStringTimeStart(),
			appointment.Doctor(),
			data.address,
			data.phone,
		)
	case RemindAboutAppointment:
		// передаём patientPhone как 6-й аргумент
		return fmt.Sprintf(
			InClinicVisitReminderTemplate,
			appointment.PatientName(),        // %s — имя пациента
			appointment.GetStringDateStart(), // %s — дата
			appointment.GetStringTimeStart(), // %s — время
			appointment.Doctor(),             // %s — врач
			data.address,                     // %s — адрес
			patientPhone,                     // %s — токен/телефон для URL
			data.phone,                       // %s — телефон клиники
		)
	}

	return ""
}

func (s *Service) getPatientPhone(ctx context.Context, appointment appointment.Appointment) (string, error) {
	patientPhone := appointment.PatientPhone()
	if len(patientPhone) == 0 {
		patientDTO, err := s.misGateway.GetPatientById(ctx, appointment.PatientId())
		if err != nil {
			return "", err
		}
		patientPhone = patientDTO.Phone
	}
	return s.cleanPhoneNumber(patientPhone), nil
}

func (s *Service) cleanPhoneNumber(phone string) string {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	return phone
}

func (s *Service) prepareMessage(template string, appointment appointment.Appointment, clinic clinicData) string {
	return fmt.Sprintf(
		template,
		appointment.GetStringDateTimeStart(),
		appointment.PatientName(),
		appointment.GetStringDateStart(),
		appointment.GetStringTimeStart(),
		appointment.Clinic(),
		appointment.Doctor(),
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

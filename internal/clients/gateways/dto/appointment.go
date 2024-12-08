package dto

import "sorkin_bot/internal/domain/entity/appointment"

type AppointmentDTO struct {
	Id               int
	TimeStart        string
	TimeEnd          string
	ClinicId         int
	Clinic           string
	DoctorId         int
	Doctor           string
	PatientId        int
	PatientName      string
	PatientBirthDate string
	PatientGender    string
	PatientPhone     string
	PatientEmail     string
	DateCreated      string
	DateUpdated      string
	Status           string
	StatusId         int
	ConfirmStatus    string
	Source           string
	MovedTo          int
	MovedFrom        int
}

func NewAppointmentDTO(id, clinicId, doctorId, patientId, statusId, movedTo, movedFrom int,
	timeStart, timeEnd, clinic, doctor, patientName, patientBirthDate, patientGender,
	patientPhone, patientEmail, dateCreated, dateUpdated, status, confirmStatus, source string,
) AppointmentDTO {
	return AppointmentDTO{
		Id:               id,
		Clinic:           clinic,
		Doctor:           doctor,
		DateUpdated:      dateUpdated,
		ClinicId:         clinicId,
		DoctorId:         doctorId,
		PatientId:        patientId,
		Status:           status,
		StatusId:         statusId,
		TimeEnd:          timeEnd,
		TimeStart:        timeStart,
		PatientEmail:     patientEmail,
		PatientName:      patientName,
		PatientPhone:     patientPhone,
		PatientGender:    patientGender,
		PatientBirthDate: patientBirthDate,
		DateCreated:      dateCreated,
		ConfirmStatus:    confirmStatus,
		Source:           source,
		MovedTo:          movedTo,
		MovedFrom:        movedFrom,
	}
}

func (a AppointmentDTO) ToDomain() appointment.Appointment {
	return appointment.NewAppointment(
		a.Id, a.ClinicId, a.DoctorId, a.PatientId, a.StatusId, a.MovedTo, a.MovedFrom, a.TimeStart, a.TimeEnd, a.Clinic, a.Doctor,
		a.PatientName, a.PatientBirthDate, a.PatientGender, a.PatientPhone, a.PatientEmail, a.DateCreated, a.DateUpdated,
		a.Status, a.ConfirmStatus, a.Source,
	)
}

type CreateAppointmentDTO struct {
	PatientId         int
	DoctorId          int
	TimeStart         string
	TimeEnd           string
	HomeAddress       string
	HomeVisit         bool
	OnlineAppointment bool
	ClinicAppointment bool
}

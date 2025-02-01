package webhook_events

import "sorkin_bot/internal/domain/entity/appointment"

type AppointmentDTO struct {
	Id               int    `json:"id"`
	TimeStart        string `json:"time_start"` //dd.mm.yyyy hh:mm
	TimeEnd          string `json:"time_end"`   //dd.mm.yyyy hh:mm
	ClinicId         int    `json:"clinic_id"`
	Clinic           string `json:"clinic"`
	DoctorId         int    `json:"doctor_id"`
	Doctor           string `json:"doctor"`
	PatientId        int    `json:"patient_id"`
	PatientName      string `json:"patient_name"`
	PatientBirthDate string `json:"patient_birth_date"`
	PatientGender    string `json:"patient_gender"`
	PatientPhone     string `json:"patient_phone"`
	DateCreated      string `json:"date_created"` //dd.mm.yyyy hh:mm
	DateUpdated      string `json:"date_updated"` //dd.mm.yyyy hh:mm
	Status           string `json:"status"`
	StatusId         int    `json:"status_id"`
	MovedFrom        int    `json:"moved_from"`
	MovedTo          int    `json:"moved_to"`
}

type AppointmentRequest struct {
	Event string         `json:"event"`
	Data  AppointmentDTO `json:"data"`
	Date  string         `json:"date"`
}

func (a *AppointmentDTO) ToDomain() appointment.Appointment {
	return appointment.NewAppointment(
		a.Id, a.ClinicId, a.DoctorId, a.PatientId, a.StatusId, a.MovedTo, a.MovedFrom, 0,
		a.TimeStart, a.TimeEnd, a.Clinic, a.Doctor, a.PatientName, a.PatientBirthDate, a.PatientGender,
		a.PatientPhone, "", a.DateCreated, a.DateUpdated, a.Status, "",
	)
}

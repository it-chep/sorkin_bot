package appointment

import "strings"

type Appointment struct {
	id               int
	timeStart        string
	timeEnd          string
	clinicId         int
	clinic           string
	doctorId         int
	doctor           string
	patientId        int
	patientName      string
	patientBirthDate string
	patientGender    string
	patientPhone     string
	patientEmail     string
	dateCreated      string
	dateUpdated      string
	status           string
	statusId         int
	confirmStatus    string
	source           string
	movedTo          int
	movedFrom        int
}

func NewAppointment(id, clinicId, doctorId, patientId, statusId, movedTo, movedFrom int,
	timeStart, timeEnd, clinic, doctor, patientName, patientBirthDate, patientGender,
	patientPhone, patientEmail, dateCreated, dateUpdated, status, confirmStatus, source string,
) Appointment {
	return Appointment{
		id:               id,
		clinic:           clinic,
		doctor:           doctor,
		dateUpdated:      dateUpdated,
		clinicId:         clinicId,
		doctorId:         doctorId,
		patientId:        patientId,
		status:           status,
		statusId:         statusId,
		timeEnd:          timeEnd,
		timeStart:        timeStart,
		patientEmail:     patientEmail,
		patientName:      patientName,
		patientPhone:     patientPhone,
		patientGender:    patientGender,
		patientBirthDate: patientBirthDate,
		dateCreated:      dateCreated,
		confirmStatus:    confirmStatus,
		source:           source,
		movedTo:          movedTo,
		movedFrom:        movedFrom,
	}
}

func (a Appointment) GetAppointmentId() int {
	return a.id
}

func (a Appointment) GetDoctorId() int {
	return a.doctorId
}

func (a Appointment) GetDate() string {
	date := strings.Split(a.timeEnd, " ")[0]
	return date
}

func (a Appointment) GetTimeStart() string {
	return a.timeStart
}

func (a Appointment) GetTimeStartShort() string {
	timeShort := strings.Split(a.timeStart, " ")[1]
	return timeShort
}

func (a Appointment) GetTimeEnd() string {
	return a.timeEnd
}

func (a Appointment) GetTimeEndShort() string {
	timeShort := strings.Split(a.timeEnd, " ")[1]
	return timeShort
}

type DraftAppointment struct {
	timeStart       *string
	timeEnd         *string
	doctorId        *int
	doctorName      *string
	tgId            *int64
	specialityId    *int
	date            *string
	appointmentType *AppointmentType
}

func NewDraftAppointment(specialityId, doctorId *int, tgId *int64, doctorName, timeStart, timeEnd, date *string, appointmentType *AppointmentType,
) DraftAppointment {
	return DraftAppointment{
		doctorId:        doctorId,
		tgId:            tgId,
		doctorName:      doctorName,
		timeEnd:         timeEnd,
		timeStart:       timeStart,
		specialityId:    specialityId,
		date:            date,
		appointmentType: appointmentType,
	}
}

func (a DraftAppointment) GetTimeStart() *string {
	return a.timeStart
}

func (a DraftAppointment) GetDoctorName() *string { return a.doctorName }

func (a DraftAppointment) GetAppointmentType() *AppointmentType {
	return a.appointmentType
}

func (a DraftAppointment) GetTimeEnd() *string {
	return a.timeEnd
}

func (a DraftAppointment) GetDoctorId() *int {
	return a.doctorId
}

func (a DraftAppointment) GetTgId() *int64 {
	return a.tgId
}

func (a DraftAppointment) GetSpecialityId() *int {
	return a.specialityId
}

func (a DraftAppointment) GetDate() *string {
	return a.date
}

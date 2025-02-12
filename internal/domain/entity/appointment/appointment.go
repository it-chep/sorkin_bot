package appointment

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Status = int

const (
	ActiveAppointment     Status = 1
	RescheduleAppointment Status = 2
	CancelAppointment     Status = 5
)

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
	confirmStatus    int
	source           string
	isOutside        int
	isTelemedicine   int
	movedTo          int
	movedFrom        int
}

func NewAppointment(id, clinicId, doctorId, patientId, statusId, movedTo, movedFrom, confirmStatus, isOutside, isTelemedicine int,
	timeStart, timeEnd, clinic, doctor, patientName, patientBirthDate, patientGender,
	patientPhone, patientEmail, dateCreated, dateUpdated, status, source string,
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
		isOutside:        isOutside,
		isTelemedicine:   isTelemedicine,
	}
}

func (a Appointment) Id() int {
	return a.id
}

func (a Appointment) MovedToID() int {
	return a.movedTo
}

func (a Appointment) MovedFromID() int {
	return a.movedFrom
}

func (a Appointment) DoctorId() int {
	return a.doctorId
}

func (a Appointment) IsOutside() bool {
	return a.isOutside != 0
}

func (a Appointment) IsTelemedicine() bool {
	return a.isTelemedicine != 0
}

func (a Appointment) GetDate() string {
	date := strings.Split(a.timeEnd, " ")[0]
	return date
}

func (a Appointment) IsCanceled() bool {
	if a.statusId == CancelAppointment {
		return true
	}

	return false
}

func (a Appointment) ClinicId() int {
	return a.clinicId
}

func (a Appointment) Doctor() string {
	return a.doctor
}

func (a Appointment) PatientName() string {
	return a.patientName
}

func (a Appointment) Clinic() string {
	return a.clinic
}

func (a Appointment) PatientId() int {
	return a.patientId
}

func (a Appointment) GetStringDateTimeStart() string {
	return a.timeStart
}

func (a Appointment) GetStringDateStart() string {
	parts := strings.Split(a.timeStart, " ")
	return parts[0]
}

func (a Appointment) GetStringTimeStart() string {
	parts := strings.Split(a.timeStart, " ")
	return parts[1]
}

func (a Appointment) GetDateTimeStart() (time.Time, error) {
	location, err := time.LoadLocation("Europe/Lisbon")
	if err != nil {
		return time.Time{}, err
	}

	parts := strings.Split(a.timeStart, " ")
	dateParts := strings.Split(parts[0], ".")
	day, err := strconv.Atoi(dateParts[0])

	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day: %w", err)
	}

	month, err := strconv.Atoi(dateParts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid month: %w", err)
	}

	year, err := strconv.Atoi(dateParts[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year: %w", err)
	}

	timeParts := strings.Split(parts[1], ":")

	hour, err := strconv.Atoi(timeParts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid hour: %w", err)
	}

	minute, err := strconv.Atoi(timeParts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid minute: %w", err)
	}

	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, location), nil
}

func (a Appointment) PatientPhone() string {
	return a.patientPhone
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

package appointment

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

func (a Appointment) GetTimeStart() string {
	return a.timeStart
}

func (a Appointment) GetTimeEnd() string {
	return a.timeEnd
}

type DraftAppointment struct {
	id           int
	timeStart    *string
	timeEnd      *string
	doctorId     *int
	tgId         *int64
	specialityId *int
}

func NewDraftAppointment(id int, specialityId, doctorId *int, tgId *int64, timeStart, timeEnd *string,
) DraftAppointment {
	return DraftAppointment{
		id:           id,
		doctorId:     doctorId,
		tgId:         tgId,
		timeEnd:      timeEnd,
		timeStart:    timeStart,
		specialityId: specialityId,
	}
}

func (a DraftAppointment) GetAppointmentId() int {
	return a.id
}

func (a DraftAppointment) GetTimeStart() *string {
	return a.timeStart
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

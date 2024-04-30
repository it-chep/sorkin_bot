package appointment

type Doctor struct {
	id                     int
	name                   string
	phone                  string
	email                  string
	professionTitles       string
	secondProfessionTitles string
	isDeleted              bool
}

func NewDoctor(id int, name, phone, email, professionTitles, secondProfessionTitles string, isDeleted bool) Doctor {
	return Doctor{
		id:                     id,
		name:                   name,
		phone:                  phone,
		email:                  email,
		professionTitles:       professionTitles,
		secondProfessionTitles: secondProfessionTitles,
		isDeleted:              isDeleted,
	}
}

type Speciality struct {
	id         int
	name       string
	doctorName string
	isDeleted  bool
}

func NewSpeciality(id int, name, doctorName string, isDeleted bool) Speciality {
	return Speciality{
		id:         id,
		name:       name,
		doctorName: doctorName,
		isDeleted:  isDeleted,
	}
}

func (s Speciality) GetDoctorName() string {
	return s.doctorName
}

func (s Speciality) GetId() int {
	return s.id
}

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
	movedTo          string
	movedFrom        string
}

func NewAppointment(id, clinicId, doctorId, patientId, statusId int,
	timeStart, timeEnd, clinic, doctor, patientName, patientBirthDate, patientGender,
	patientPhone, patientEmail, dateCreated, dateUpdated, status, confirmStatus, source, movedTo, movedFrom string,
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

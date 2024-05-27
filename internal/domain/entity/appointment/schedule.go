package appointment

type Schedule struct {
	clinicId       int
	doctorId       int
	date           string
	timeStart      string
	timeStartShort string
	timeEnd        string
	timeEndShort   string
	category       string
	user           string
	categoryId     int
	profession     string
	room           string
	isBusy         bool
	isPast         bool
}

func NewSchedule(clinicId, doctorId, categoryId int,
	date, timeStart, timeStartShort, timeEnd, timeEndShort, category, profession, room, user string,
	isBusy, isPast bool,
) Schedule {
	return Schedule{
		clinicId:       clinicId,
		doctorId:       doctorId,
		date:           date,
		timeStart:      timeStart,
		timeStartShort: timeEndShort,
		timeEnd:        timeEnd,
		timeEndShort:   timeStartShort,
		category:       category,
		categoryId:     categoryId,
		profession:     profession,
		room:           room,
		isBusy:         isBusy,
		isPast:         isPast,
		user:           user,
	}
}

func (sch Schedule) GetProfession() string {
	return sch.profession
}

func (sch Schedule) GetTimeStart() string {
	return sch.timeStart
}

func (sch Schedule) GetTimeStartShort() string {
	return sch.timeStartShort
}

func (sch Schedule) GetTimeEnd() string {
	return sch.timeEnd
}

func (sch Schedule) GetTimeEndShort() string {
	return sch.timeEndShort
}

func (sch Schedule) GetDoctorId() int {
	return sch.doctorId
}

func (sch Schedule) GetDate() string {
	return sch.date
}

func (sch Schedule) GetDoctorName() string {
	return sch.user
}

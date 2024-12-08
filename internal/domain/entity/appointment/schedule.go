package appointment

import (
	"fmt"
	"strings"
	"time"
)

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

type SchedulePeriod struct {
	date      string
	doctorId  int
	timeStart string
	timeEnd   string
}

func NewSchedulePeriod(date, timeStart, timeEnd string, doctorId int) SchedulePeriod {
	return SchedulePeriod{
		date:      date,
		doctorId:  doctorId,
		timeStart: timeStart,
		timeEnd:   timeEnd,
	}
}

func (sp *SchedulePeriod) GetDate() string {
	return sp.date
}

func (sp *SchedulePeriod) GetDoctorId() int {
	return sp.doctorId
}

func (sp *SchedulePeriod) GetTimeStart() string {
	return sp.timeStart
}

func (sp *SchedulePeriod) GetTimeEnd() string {
	return sp.timeEnd
}

func (sp *SchedulePeriod) GetDateInTimeType() time.Time {
	dateItems := strings.Split(sp.date, ".")
	dateStr := fmt.Sprintf("%s-%s-%s", dateItems[2], dateItems[1], dateItems[0])

	parsedDate, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return time.Time{}
	}
	return parsedDate
}

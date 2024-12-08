package dto

import (
	"sorkin_bot/internal/domain/entity/appointment"
	"strconv"
	"strings"
)

type ScheduleDTO struct {
	clinicId       int
	doctorId       int
	date           string
	timeStart      string
	timeStartShort string
	timeEnd        string
	timeEndShort   string
	category       string
	categoryId     int
	profession     string
	room           string
	user           string
	isBusy         bool
	isPast         bool
}

func NewScheduleDTO(clinicId, doctorId, categoryId int,
	date, timeStart, timeStartShort, timeEnd, timeEndShort, category, profession, room, user string,
	isBusy, isPast bool,
) ScheduleDTO {
	return ScheduleDTO{
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
		user:           user,
		isBusy:         isBusy,
		isPast:         isPast,
	}
}

func (d ScheduleDTO) ToDomain() appointment.Schedule {
	return appointment.NewSchedule(
		d.clinicId, d.doctorId, d.categoryId, d.date, d.timeStart, d.timeStartShort, d.timeEnd,
		d.timeEndShort, d.category, d.profession, d.room, d.user, d.isBusy, d.isPast,
	)
}

func (d ScheduleDTO) GetTimeStart() string {
	return d.timeStart
}

func (d ScheduleDTO) GetTimeStartShort() string {
	return d.timeStartShort
}

func (d ScheduleDTO) GetIntHourStart() int {
	items := strings.Split(d.timeStartShort, ":")
	hour, _ := strconv.Atoi(items[0])
	return hour
}

type SchedulePeriodDTO struct {
	date      string
	doctorId  int
	timeStart string
	timeEnd   string
}

func NewSchedulePeriodDTO(date, timeStart, timeEnd string, doctorId int) SchedulePeriodDTO {
	return SchedulePeriodDTO{
		date:      date,
		timeStart: timeStart,
		timeEnd:   timeEnd,
		doctorId:  doctorId,
	}
}

func (d SchedulePeriodDTO) ToDomain() appointment.SchedulePeriod {
	return appointment.NewSchedulePeriod(
		d.date, d.timeStart, d.timeEnd, d.doctorId,
	)
}

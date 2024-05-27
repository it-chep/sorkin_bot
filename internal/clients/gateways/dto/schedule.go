package dto

import "sorkin_bot/internal/domain/entity/appointment"

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

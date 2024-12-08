package appointment

import (
	"context"
	"fmt"
	"math/rand"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"strings"
	"time"
)

func (as *AppointmentService) GetFastAppointmentSchedules(ctx context.Context) (randomDoctors map[int]appointment.Schedule) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetFastAppointmentSchedules"

	currentTime := time.Now()
	timeStart := fmt.Sprintf("%02d.%02d.%d %02d:%02d", currentTime.Day()+1, currentTime.Month(), currentTime.Year(), currentTime.Hour(), currentTime.Minute())
	timeEnd := fmt.Sprintf("%02d.%02d.%d %02d:%02d", currentTime.Day()+1, currentTime.Month(), currentTime.Year(), currentTime.Hour(), currentTime.Minute())
	randomDoctors = make(map[int]appointment.Schedule)
	filteredDoctors := make(map[int][]appointment.Schedule)

	schedulesMap, err := as.misAdapter.GetSchedulesByDoctorId(ctx, 0, timeStart, timeEnd)

	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return
	}

	var doctorIDs []int
	for doctorID := range schedulesMap {
		doctorIDs = append(doctorIDs, doctorID)
	}

	// TODO УЖАС !!!!! А не вложенность, подумать, но такие кривые данные приходят от МИС
	for _, schedule := range schedulesMap {
		for _, scheduleItem := range schedule {
			for _, doctorID := range doctorIDs {
				if scheduleItem.GetDoctorId() == doctorID {
					filteredDoctors[doctorID] = append(filteredDoctors[doctorID], scheduleItem)
				}
			}
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(doctorIDs), func(i, j int) {
		doctorIDs[i], doctorIDs[j] = doctorIDs[j], doctorIDs[i]
	})

	numDoctorsToSelect := 6
	if len(doctorIDs) < 6 {
		numDoctorsToSelect = len(doctorIDs)
	}

	for i := 0; i < numDoctorsToSelect; i++ {
		doctorID := doctorIDs[i]
		for _, schedule := range filteredDoctors[doctorID] {
			randomDoctors[doctorID] = schedule
			break
		}
	}

	return randomDoctors
}

func (as *AppointmentService) GetSchedulesByDoctorId(ctx context.Context, userEntity entity.User, dayStart time.Time, doctorId *int) (schedulesMap []appointment.Schedule, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetSchedulesByDoctorId"
	timeStart := fmt.Sprintf("%02d.%02d.%d 00:00", dayStart.Day(), int(dayStart.Month()), dayStart.Year())
	timeEnd := fmt.Sprintf("%02d.%02d.%d 23:59", dayStart.Day(), int(dayStart.Month()), dayStart.Year())

	if doctorId == nil {
		draftAppointment, err := as.GetDraftAppointment(ctx, userEntity.GetTgId())
		if err != nil {
			return nil, err
		}
		doctorIdValue := draftAppointment.GetDoctorId()
		doctorId = doctorIdValue
	}

	schedules, err := as.misAdapter.GetSchedulesByDoctorId(ctx, *doctorId, timeStart, timeEnd)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return nil, err
	}

	return schedules[*doctorId], err
}

func (as *AppointmentService) GetSchedulePeriodsByDoctorId(ctx context.Context, doctorId int, dayStart time.Time) (schedulePeriodsMap map[time.Time]bool, err error) {
	timeStart := fmt.Sprintf("%02d.%02d.%d 00:00", dayStart.Day(), int(dayStart.Month()), dayStart.Year())
	nextMonth := int(dayStart.Month()) + 1
	year := dayStart.Year()
	if nextMonth == 13 {
		nextMonth = 1
		year = year + 1
	}
	timeEnd := fmt.Sprintf("01.%02d.%d 00:00", nextMonth, year)
	// мапа вида {"день месяца": работает врач или нет}
	schedulePeriodsMap = make(map[time.Time]bool, 31)

	schedulePeriods, err := as.misAdapter.GetSchedulePeriodsByDoctorId(ctx, doctorId, timeStart, timeEnd)
	if err != nil {
		return nil, err
	}

	for _, schedulePeriod := range schedulePeriods {
		schedulePeriodsMap[schedulePeriod.GetDateInTimeType()] = true
	}

	return schedulePeriodsMap, nil
}

func (as *AppointmentService) getSchedulesToHomeVisit(ctx context.Context, doctorIds []int, dayStart time.Time) (schedulesMap []appointment.Schedule, err error) {
	timeStart := fmt.Sprintf("%02d.%02d.%d 00:00", dayStart.Day(), int(dayStart.Month()), dayStart.Year())
	timeEnd := fmt.Sprintf("%02d.%02d.%d 23:59", dayStart.Day(), int(dayStart.Month()), dayStart.Year())

	availableDoctorIds, err := as.misAdapter.GetAvailableDoctorIdsFromSchedulePeriods(ctx, doctorIds, timeStart, timeEnd)
	if err != nil {
		return
	}
	schedules, err := as.misAdapter.GetSchedulesManyDoctors(ctx, availableDoctorIds, timeStart, timeEnd)
	if err != nil {
		return nil, err
	}

	return schedules, err
}

func (as *AppointmentService) GetSchedulesToHomeVisit(ctx context.Context, userEntity entity.User, dayStart time.Time) (schedulesMap []appointment.Schedule, err error) {
	doctors, err := as.misAdapter.GetDoctors(ctx, true, false, false)
	if err != nil {
		return nil, err
	}

	ids := make([]int, 0, len(doctors))
	for _, doctor := range doctors {
		if (*userEntity.GetState() == "pediatrician" && strings.Contains(doctor.GetSecondProfessionTitles(), "детское здоровье")) ||
			(*userEntity.GetState() == "therapist" && !strings.Contains(doctor.GetSecondProfessionTitles(), "детское здоровье")) {
			ids = append(ids, doctor.GetID())
		}
	}

	return as.getSchedulesToHomeVisit(ctx, ids, dayStart)
}

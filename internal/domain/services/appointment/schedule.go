package appointment

import (
	"context"
	"fmt"
	"math/rand"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
	"time"
)

func (as *AppointmentService) GetFastAppointmentSchedules(ctx context.Context) (randomDoctors map[int]appointment.Schedule) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetFastAppointmentSchedules"

	currentTime := time.Now()
	timeStart := fmt.Sprintf("%02d.%02d.%d %02d:%02d", currentTime.Day()+1, currentTime.Month(), currentTime.Year(), currentTime.Hour(), currentTime.Minute())

	randomDoctors = make(map[int]appointment.Schedule)
	filteredDoctors := make(map[int][]appointment.Schedule)

	schedulesMap, err := as.misAdapter.GetSchedules(ctx, 0, timeStart)

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

func (as *AppointmentService) GetSchedules(ctx context.Context, userEntity entity.User, doctorId *int) (schedulesMap []appointment.Schedule, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetSchedules	"

	if doctorId == nil {
		draftAppointment, err := as.GetDraftAppointment(ctx, userEntity.GetTgId())
		if err != nil {
			return nil, err
		}
		doctorIdValue := draftAppointment.GetDoctorId()
		doctorId = doctorIdValue
	}

	schedules, err := as.misAdapter.GetSchedules(ctx, *doctorId, "")
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return nil, err
	}

	return schedules[*doctorId], err

}

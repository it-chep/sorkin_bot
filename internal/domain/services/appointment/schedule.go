package appointment

import (
	"context"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
)

func (as *AppointmentService) GetFastAppointmentSchedules(ctx context.Context) (schedulesMap map[int][]appointment.Schedule) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetFastAppointmentSchedules"
	schedulesMap, err := as.misAdapter.FastAppointment(ctx)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return
	}

	return schedulesMap
}

func (as *AppointmentService) GetSchedules(ctx context.Context, doctorId int) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetSchedules	"
	schedules, err := as.misAdapter.GetSchedules(ctx, doctorId, "")
	// +- такая логика
	if doctorId == 0 {
		for doctorId, doctorSchedules := range schedules {
			as.logger.Info(fmt.Sprintf("schedules[doctorId] %d %s", doctorId, doctorSchedules))
		}
	} else {
		as.logger.Info(fmt.Sprintf("schedules[doctorId] %d %s", doctorId, schedules[doctorId]))
	}

	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return
	}
}

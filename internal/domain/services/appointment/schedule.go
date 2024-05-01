package appointment

import (
	"context"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
)

func (as AppointmentService) GetFastAppointmentSchedules(ctx context.Context) (schedulesMap map[int][]appointment.Schedule) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetFastAppointmentSchedules"
	err, schedulesMap := as.mis.FastAppointment(ctx)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return
	}

	return schedulesMap
}

func (as AppointmentService) GetSchedules(ctx context.Context, doctorId int) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetSchedules	"
	err, schedules := as.mis.GetSchedules(ctx, doctorId, "")
	as.logger.Info(fmt.Sprintf("schedules[doctorId] %s", schedules[doctorId]))

	if err != nil {
		as.logger.Error(fmt.Sprintf("error: %s. Place %s", err, op))
		return
	}
}

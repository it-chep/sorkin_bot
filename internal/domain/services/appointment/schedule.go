package appointment

import (
	"context"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
	entity "sorkin_bot/internal/domain/entity/user"
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

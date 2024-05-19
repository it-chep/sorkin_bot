package adapter

import (
	"context"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
	"time"
)

func (a *AppointmentServiceAdapter) GetSchedules(ctx context.Context, doctorId int, timeStart string) (schedulesMap map[int][]appointment.Schedule, err error) {
	cachedSchedules, ok := a.cache.Get(fmt.Sprintf("%d_schedules", doctorId))
	if !ok {
		schedulesMap = make(map[int][]appointment.Schedule)
		schedulesDTO, err := a.gateway.GetSchedules(ctx, doctorId, timeStart)
		if err != nil {
			return nil, err
		}
		for i, scheduleDTO := range schedulesDTO {
			schedules := make([]appointment.Schedule, 0, len(scheduleDTO))
			for _, scheduleItem := range scheduleDTO {
				scheduleEntity := scheduleItem.ToDomain()
				schedules = append(schedules, scheduleEntity)
			}
			schedulesMap[i] = schedules
		}
		a.cache.Set(fmt.Sprintf("%d_schedules", doctorId), schedulesMap, 10*time.Minute)
		return schedulesMap, nil
	}
	return cachedSchedules.(map[int][]appointment.Schedule), err
}

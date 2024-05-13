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
		schedulesDTO, err := a.gateway.GetSchedules(ctx, doctorId, timeStart)
		if err != nil {
			return nil, err
		}
		for i, scheduleDTO := range schedulesDTO {
			fmt.Println(i, scheduleDTO)
			//schedulesMap[i] = []scheduleDTO.ToDomain()
		}
		a.cache.Set(fmt.Sprintf("%d_schedules", doctorId), schedulesMap, 10*time.Minute)
		return schedulesMap, nil
	}
	return cachedSchedules.(map[int][]appointment.Schedule), err
}

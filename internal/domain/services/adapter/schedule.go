package adapter

import (
	"context"
	"fmt"
	"sorkin_bot/internal/domain/entity/appointment"
	"time"
)

func (a *AppointmentServiceAdapter) GetSchedulesByDoctorId(ctx context.Context, doctorId int, timeStart, timeEnd string) (schedulesMap map[int][]appointment.Schedule, err error) {
	cachedSchedules, ok := a.cache.Get(fmt.Sprintf("%d_schedules", doctorId))
	if !ok {
		schedulesMap = make(map[int][]appointment.Schedule)
		schedulesDTO, err := a.gateway.GetSchedulesByDoctorId(ctx, doctorId, timeStart, timeEnd)
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

func (a *AppointmentServiceAdapter) GetSchedulesManyDoctors(ctx context.Context, doctorIds []int, timeStart, timeEnd string) (schedules []appointment.Schedule, err error) {
	schedulesDTO, err := a.gateway.GetSchedulesManyDoctors(ctx, doctorIds, timeStart, timeEnd)
	schedules = make([]appointment.Schedule, 0, len(schedulesDTO))
	for _, scheduleItem := range schedulesDTO {
		scheduleEntity := scheduleItem.ToDomain()
		schedules = append(schedules, scheduleEntity)
	}
	return schedules, nil
}

func (a *AppointmentServiceAdapter) GetAvailableDoctorIdsFromSchedulePeriods(ctx context.Context, doctorIds []int, timeStart, timeEnd string) (availableIds []int, err error) {
	// todo cache
	return a.gateway.GetAvailableDoctorIdsFromSchedulePeriods(ctx, doctorIds, timeStart, timeEnd)
}

func (a *AppointmentServiceAdapter) GetSchedulePeriodsByDoctorId(ctx context.Context, doctorId int, timeStart, timeEnd string) (schedulePeriods []appointment.SchedulePeriod, err error) {
	// todo cache
	schedulePeriodsDTO, err := a.gateway.GetSchedulePeriodsByDoctorId(ctx, doctorId, timeStart, timeEnd)
	if err != nil {
		return nil, err
	}
	for _, schedulePeriodDTO := range schedulePeriodsDTO {
		schedulePeriods = append(schedulePeriods, schedulePeriodDTO.ToDomain())
	}
	return schedulePeriods, nil
}

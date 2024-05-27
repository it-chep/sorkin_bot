package mis_reno

import (
	"context"
	"sorkin_bot/internal/clients/gateways/dto"
	"sorkin_bot/internal/clients/gateways/mis_reno/mis_dto"
	"strconv"
)

func (mg *MisRenoGateway) GetSchedules(ctx context.Context, doctorId int, timeStart string) (schedulesMap map[int][]dto.ScheduleDTO, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetSchedules"
	var (
		request = mis_dto.GetScheduleRequest{
			DoctorId:   doctorId,
			ShowBusy:   false,
			ShowPast:   false,
			ShowAll:    false,
			AllClinics: false,
		}
		response  mis_dto.GetScheduleResponse
		schedules []dto.ScheduleDTO
	)

	// Если мы хотим получить по конкретному дню, то добавляем параметр timeStart
	if timeStart != "" {
		request.TimeStart = timeStart
	}

	schedulesMap = make(map[int][]dto.ScheduleDTO)

	responseBody := mg.sendToMIS(ctx, mis_dto.GetScheduleMethod, JsonMarshaller(request, op, mg.logger))

	response, err = JsonUnMarshaller(response, responseBody, op, mg.logger)
	if err != nil {
		return schedulesMap, err
	}

	// Если мы берем расписание у всех докторов
	for strResponseDoctorId, doctorSchedule := range response.Data {
		for _, schedule := range doctorSchedule {
			schedules = append(schedules, schedule.ToDTO())
		}
		if doctorId == 0 {
			responseDoctorId, _ := strconv.Atoi(strResponseDoctorId)
			schedulesMap[responseDoctorId] = schedules
		} else {
			schedulesMap[doctorId] = schedules
		}
	}

	return schedulesMap, nil
}

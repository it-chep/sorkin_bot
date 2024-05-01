package mis_reno

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	mis_dto "sorkin_bot/internal/clients/gateways/mis_reno/mis_dto"
	"sorkin_bot/internal/domain/entity/appointment"
	"strconv"
)

func (mg *MisRenoGateway) GetSchedules(ctx context.Context, doctorId int, timeStart string) (err error, schedulesMap map[int][]appointment.Schedule) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetSchedules"
	var response mis_dto.GetScheduleResponse
	var schedules []appointment.Schedule
	var request = mis_dto.GetScheduleRequest{
		DoctorId:   doctorId,
		ShowBusy:   false,
		ShowPast:   false,
		ShowAll:    false,
		AllClinics: false,
	}

	// Если мы хотим получить по конкретному дню, то добавляем параметр timeStart
	if timeStart != "" {
		request.TimeStart = timeStart
	}

	schedulesMap = make(map[int][]appointment.Schedule)

	jsonBody, err := json.Marshal(request)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err, schedulesMap
	}
	body := bytes.NewReader(jsonBody)
	responseBody := mg.sendToMIS(ctx, mis_dto.GetScheduleMethod, body)

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Info(fmt.Sprintf("error while unmarshalling json %s \nplace: %s", err, op))
		return err, schedulesMap
	}

	if doctorId != 0 {
		// Если мы берем расписание у конкретного доктора
		for _, schedule := range response.Data[fmt.Sprintf("%d", doctorId)] {
			schedules = append(schedules, schedule.ToDomain())
		}
		schedulesMap[doctorId] = schedules
	} else {
		// Если мы берем расписание у всех докторов
		for strResponseDoctorId, doctorSchedule := range response.Data {
			for _, schedule := range doctorSchedule {
				schedules = append(schedules, schedule.ToDomain())
			}
			responseDoctorId, _ := strconv.Atoi(strResponseDoctorId)
			schedulesMap[responseDoctorId] = schedules
		}
	}

	return nil, schedulesMap
}

func (mg *MisRenoGateway) ChooseSchedules(ctx context.Context) {

}

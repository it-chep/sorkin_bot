package mis_reno

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	mis_dto "sorkin_bot/internal/clients/gateways/mis_reno/mis_dto"
	"sorkin_bot/internal/domain/entity/appointment"
)

func (mg *MisRenoGateway) GetSchedules(ctx context.Context, doctorId int) (err error, schedules []appointment.Schedule) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetSchedules"
	var response mis_dto.GetScheduleResponse
	var request = mis_dto.GetScheduleRequest{
		DoctorId:   doctorId,
		ShowBusy:   false,
		ShowPast:   false,
		ShowAll:    false,
		AllClinics: false,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err, schedules
	}
	body := bytes.NewReader(jsonBody)
	responseBody := mg.sendToMIS(ctx, mis_dto.GetScheduleMethod, body)

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Info(fmt.Sprintf("error while unmarshalling json %s \nplace: %s", err, op))
		return err, schedules
	}

	if doctorId != -1 {
		// Если мы берем расписание у конкретного доктора
		for _, schedule := range response.Data[fmt.Sprintf("%d", doctorId)] {
			schedules = append(schedules, schedule.ToDomain())
		}
	} else {
		// Если мы берем расписание у всех докторов
		for _, doctorSchedule := range response.Data {
			for _, schedule := range doctorSchedule {
				schedules = append(schedules, schedule.ToDomain())
			}
		}
	}

	return nil, schedules
}

func (mg *MisRenoGateway) ChooseSchedules(ctx context.Context) {

}

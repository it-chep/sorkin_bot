package appointment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sorkin_bot/internal/domain/services/appointment/mis_dto"
)

func (as *AppointmentService) GetSchedules(ctx context.Context) (err error) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetSchedules"
	var Request mis_dto.GetScheduleRequest
	var Response mis_dto.GetSpecialityResponse

	Request = mis_dto.GetScheduleRequest{}
	jsonBody, err := json.Marshal(Request)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err
	}
	body := bytes.NewReader(jsonBody)
	responseBody := as.sendToMIS(ctx, mis_dto.GetScheduleMethod, body)

	err = json.Unmarshal(responseBody, &Response)
	if err != nil {
		as.logger.Info(fmt.Sprintf("error while unmarshalling json %s \nplace: %s", err, op))
		return err
	}

	for _, schedule := range Response.Data {
		fmt.Printf("ID: %d\n", schedule.ID)
		fmt.Printf("Расписание: %s\n", schedule.DoctorName)
	}
	return nil
}

func (as *AppointmentService) ChooseSchedules(ctx context.Context) {

}

package mis_reno

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	mis_dto "sorkin_bot/internal/clients/gateways/mis_reno/mis_dto"
)

func (mg *MisRenoGateway) GetSchedules(ctx context.Context) (err error, response mis_dto.GetSpecialityResponse) {
	op := "sorkin_bot.internal.domain.services.appointment.schedule.GetSchedules"
	var request = mis_dto.GetScheduleRequest{}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err, response
	}
	body := bytes.NewReader(jsonBody)
	responseBody := mg.sendToMIS(ctx, mis_dto.GetScheduleMethod, body)

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Info(fmt.Sprintf("error while unmarshalling json %s \nplace: %s", err, op))
		return err, response
	}

	for _, schedule := range response.Data {
		fmt.Printf("ID: %d\n", schedule.ID)
		fmt.Printf("Расписание: %s\n", schedule.DoctorName)
	}
	return nil, response
}

func (mg *MisRenoGateway) ChooseSchedules(ctx context.Context) {

}

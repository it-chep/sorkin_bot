package mis_reno

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sorkin_bot/internal/clients/gateways/mis_reno/mis_dto"
)

func (mg *MisRenoGateway) GetPatientById(ctx context.Context) (err error) {
	op := "sorkin_bot.internal.domain.services.appointment.users.GetPatientById"
	var response mis_dto.GetPatientResponse
	var request = mis_dto.GetPatientRequest{}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err
	}
	body := bytes.NewReader(jsonBody)
	responseBody := mg.sendToMIS(ctx, mis_dto.GetPatientMethod, body)

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Info(fmt.Sprintf("error while unmarshalling json %s \nplace: %s", err, op))
		return err
	}

	return nil
}

func (mg *MisRenoGateway) CreatePatient(ctx context.Context) (err error) {
	op := "sorkin_bot.internal.domain.services.appointment.users.CreatePatient"
	var response mis_dto.CreatePatientResponse
	var request = mis_dto.CreatePatientRequest{}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err
	}
	body := bytes.NewReader(jsonBody)
	responseBody := mg.sendToMIS(ctx, mis_dto.CreatePatientMethod, body)

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Info(fmt.Sprintf("error while unmarshalling json %s \nplace: %s", err, op))
		return err
	}

	return nil
}

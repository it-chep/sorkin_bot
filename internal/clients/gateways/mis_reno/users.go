package mis_reno

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sorkin_bot/internal/clients/gateways/mis_reno/mis_dto"
	entity "sorkin_bot/internal/domain/entity/user"
)

func (mg *MisRenoGateway) GetPatientById(ctx context.Context, patientId int) (err error) {
	op := "sorkin_bot.internal.domain.services.appointment.users.GetPatientById"
	var response mis_dto.GetPatientResponse
	var request = mis_dto.GetPatientRequest{
		Id: patientId,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err
	}
	body := bytes.NewReader(jsonBody)
	responseBody := mg.sendToMIS(ctx, mis_dto.GetPatientMethod, body)

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while unmarshalling json %s \nplace: %s", err, op))
		return err
	}

	return nil
}

func (mg *MisRenoGateway) CreatePatient(ctx context.Context, user entity.User) (err error, patientId int) {
	op := "sorkin_bot.internal.domain.services.appointment.users.CreatePatient"
	var response mis_dto.CreatePatientResponse
	//todo обязательно продебажить CreatePatientRequest
	var request = mis_dto.CreatePatientRequest{
		LastName:  user.GetLastName(),
		ThirdName: user.GetThirdName(),
		FirstName: user.GetFirstName(),
		Phone:     user.GetPhone(),
		BirthDate: user.GetBirthDate(),
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err, -1
	}
	body := bytes.NewReader(jsonBody)
	responseBody := mg.sendToMIS(ctx, mis_dto.CreatePatientMethod, body)

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Info(fmt.Sprintf("error while unmarshalling json %s \nplace: %s", err, op))
		return err, -1
	}

	return nil, response.Data.PatientID
}

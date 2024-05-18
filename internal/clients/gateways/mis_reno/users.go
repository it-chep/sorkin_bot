package mis_reno

import (
	"context"
	"sorkin_bot/internal/clients/gateways/dto"
	"sorkin_bot/internal/clients/gateways/mis_reno/mis_dto"
)

func (mg *MisRenoGateway) GetPatientById(ctx context.Context, patientId int) (patientDTO dto.CreatedPatientDTO, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.users.GetPatientById"
	var response mis_dto.MisGetPatientResponse
	var request = mis_dto.GetPatientRequest{
		Id: patientId,
	}
	responseBody := mg.sendToMIS(ctx, mis_dto.GetPatientMethod, JsonMarshaller(request, op, mg.logger))

	response, err = JsonUnMarshaller(response, responseBody, op, mg.logger)
	if err != nil {
		return patientDTO, err
	}

	return response.Data.ToDTO(), nil
}

func (mg *MisRenoGateway) CreatePatient(ctx context.Context, userDTO dto.PatientDTO) (patientId *int, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.users.CreatePatient"
	var response mis_dto.CreatePatientResponse
	//todo обязательно продебажить CreatePatientRequest
	var request = mis_dto.CreatePatientRequest{
		LastName:  userDTO.LastName,
		ThirdName: userDTO.ThirdName,
		FirstName: userDTO.FirstName,
		Phone:     userDTO.Phone,
		BirthDate: userDTO.BirthDate,
	}

	responseBody := mg.sendToMIS(ctx, mis_dto.CreatePatientMethod, JsonMarshaller(request, op, mg.logger))

	response, err = JsonUnMarshaller(response, responseBody, op, mg.logger)
	if err != nil {
		return &response.Data.PatientID, err
	}

	return &response.Data.PatientID, nil
}

package mis_reno

import (
	"context"
	"sorkin_bot/internal/clients/gateways/dto"
	"sorkin_bot/internal/clients/gateways/mis_reno/mis_dto"
	entity "sorkin_bot/internal/domain/entity/user"
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
		ThirdName: ".",
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

func (mg *MisRenoGateway) GetPatientByBirthDate(ctx context.Context, user entity.User) (patientDTO dto.CreatedPatientDTO, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.users.GetPatientByBirthDate"
	var (
		responseOne  mis_dto.MisGetPatientResponse
		responseMany mis_dto.MisGetPatientsResponse
	)

	var request = mis_dto.GetPatientRequest{
		BirthDate: *user.GetBirthDate(),
	}
	responseBody := mg.sendToMIS(ctx, mis_dto.GetPatientMethod, JsonMarshaller(request, op, mg.logger))

	responseOne, err = JsonUnMarshaller(responseOne, responseBody, op, mg.logger)
	if err == nil {
		if responseOne.Data.LastName == *user.GetLastName() && responseOne.Data.FirstName == user.GetFirstName() || responseOne.Data.LastName == user.GetFirstName() && responseOne.Data.FirstName == *user.GetLastName() {
			patientDTO = responseOne.Data.ToDTO()
			return patientDTO, nil
		}
	}

	responseMany, err = JsonUnMarshaller(responseMany, responseBody, op, mg.logger)
	if err != nil {
		return patientDTO, err
	}
	for _, patient := range responseMany.Data {
		if patient.LastName == *user.GetLastName() && patient.FirstName == user.GetFirstName() || patient.LastName == user.GetFirstName() && patient.FirstName == *user.GetLastName() {
			patientDTO = patient.ToDTO()
			break
		}
	}

	return patientDTO, nil
}

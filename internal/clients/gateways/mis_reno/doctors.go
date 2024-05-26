package mis_reno

import (
	"context"
	"encoding/json"
	"fmt"
	"sorkin_bot/internal/clients/gateways/dto"
	"sorkin_bot/internal/clients/gateways/mis_reno/mis_dto"
)

func (mg *MisRenoGateway) GetDoctorsBySpecialityId(ctx context.Context, specialityId int) (doctors []dto.DoctorDTO, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.doctors.GetDoctorsBySpecialityId"
	var response mis_dto.GetUsersResponse
	var request = mis_dto.GetUserRequest{
		SpecialityId: specialityId,
	}

	responseBody := mg.sendToMIS(ctx, mis_dto.GetUsersMethod, JsonMarshaller(request, op, mg.logger))

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Info(fmt.Sprintf("error while unmarshalling json %s \nplace: %s", err, op))
		return doctors, err
	}

	for _, doctor := range response.Data {
		doctors = append(doctors, doctor.ToDTO())
	}

	return doctors, nil
}

func (mg *MisRenoGateway) GetDoctorInfo(ctx context.Context, doctorId int) (doctorDTO dto.DoctorDTO, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.doctors.GetDoctorInfo"
	var response mis_dto.GetUsersResponse
	var request = mis_dto.GetUserRequest{
		DoctorId: doctorId,
	}

	responseBody := mg.sendToMIS(ctx, mis_dto.GetUsersMethod, JsonMarshaller(request, op, mg.logger))

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Info(fmt.Sprintf("error while unmarshalling json %s \nplace: %s", err, op))
		return doctorDTO, err
	}

	for _, doctor := range response.Data {
		if doctor.MisUser.ID == doctorId {
			doctorDTO = doctor.ToDTO()
			break
		}
	}

	return doctorDTO, nil
}

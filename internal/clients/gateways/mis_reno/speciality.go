package mis_reno

import (
	"context"
	"sorkin_bot/internal/clients/gateways/dto"
	"sorkin_bot/internal/clients/gateways/mis_reno/mis_dto"
)

func (mg *MisRenoGateway) GetSpecialities(ctx context.Context) (specialities []dto.SpecialityDTO, err error) {
	op := "sorkin_bot.internal.domain.services.appointment.speciality.GetSpecialities"
	var response mis_dto.GetSpecialityResponse
	var request = mis_dto.GetSpecialityRequest{
		ShowAll:        false,
		ShowDeleted:    false,
		WithoutDoctors: false,
	}

	responseBody := mg.sendToMIS(ctx, mis_dto.GetSpecialityMethod, JsonMarshaller(request, op, mg.logger))

	response, err = JsonUnMarshaller(response, responseBody, op, mg.logger)
	if err != nil {
		return specialities, err
	}

	for _, specialityMis := range response.Data {
		specialities = append(specialities, specialityMis.ToDTO())
	}

	return specialities, nil
}

func (mg *MisRenoGateway) ChooseSpecialities(ctx context.Context) {

}

package mis_reno

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	mis_dto "sorkin_bot/internal/clients/gateways/mis_reno/mis_dto"
	"sorkin_bot/internal/domain/entity/appointment"
)

func (mg *MisRenoGateway) GetSpecialities(ctx context.Context) (err error, specialities []appointment.Speciality) {
	op := "sorkin_bot.internal.domain.services.appointment.speciality.GetSpecialities"
	var response mis_dto.GetSpecialityResponse
	var request = mis_dto.GetSpecialityRequest{
		ShowAll:        false,
		ShowDeleted:    false,
		WithoutDoctors: false,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err, specialities
	}
	body := bytes.NewReader(jsonBody)
	responseBody := mg.sendToMIS(ctx, mis_dto.GetSpecialityMethod, body)

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Info(fmt.Sprintf("error while unmarshalling json %s \nplace: %s", err, op))
		return err, specialities
	}

	for _, specialityDTO := range response.Data {
		specialities = append(specialities, specialityDTO.ToDomain())
	}

	return nil, specialities
}

func (mg *MisRenoGateway) ChooseSpecialities(ctx context.Context) {

}

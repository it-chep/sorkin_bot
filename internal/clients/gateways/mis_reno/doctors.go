package mis_reno

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	mis_dto "sorkin_bot/internal/clients/gateways/mis_reno/mis_dto"
	"sorkin_bot/internal/domain/entity/appointment"
)

func (mg *MisRenoGateway) GetDoctors(ctx context.Context, specialityId int) (err error, doctors []appointment.Doctor) {
	op := "sorkin_bot.internal.domain.services.appointment.doctors.GetDoctors"
	var response mis_dto.GetUsersResponse
	var request = mis_dto.GetUserRequest{
		SpecialityId: specialityId,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		mg.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err, doctors
	}
	body := bytes.NewReader(jsonBody)
	responseBody := mg.sendToMIS(ctx, mis_dto.GetUsersMethod, body)

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		mg.logger.Info(fmt.Sprintf("error while unmarshalling json %s \nplace: %s", err, op))
		return err, doctors
	}
	for _, doctorDTO := range response.Data {
		doctors = append(doctors, doctorDTO.ToDomain())
	}
	return nil, doctors
}

func (mg *MisRenoGateway) ChooseDoctors(ctx context.Context) {

}

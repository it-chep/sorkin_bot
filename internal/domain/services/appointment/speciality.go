package appointment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sorkin_bot/internal/domain/services/appointment/mis_dto"
)

func (as *AppointmentService) GetSpecialities(ctx context.Context) (err error) {
	op := "sorkin_bot.internal.domain.services.appointment.speciality.GetSpecialities"
	var Request mis_dto.GetSpecialityRequest
	var Response mis_dto.GetSpecialityResponse

	Request = mis_dto.GetSpecialityRequest{
		ShowAll:        false,
		ShowDeleted:    false,
		WithoutDoctors: false,
	}
	jsonBody, err := json.Marshal(Request)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err
	}
	body := bytes.NewReader(jsonBody)
	responseBody := as.sendToMIS(ctx, mis_dto.GetSpecialityMethod, body)

	err = json.Unmarshal(responseBody, &Response)
	if err != nil {
		as.logger.Info(fmt.Sprintf("error while unmarshalling json %s \nplace: %s", err, op))
		return err
	}

	for _, speciality := range Response.Data {
		fmt.Printf("ID: %d\n", speciality.ID)
		fmt.Printf("Специальность: %s\n", speciality.DoctorName)
	}
	return nil
}

func (as *AppointmentService) ChooseSpecialities(ctx context.Context) {

}

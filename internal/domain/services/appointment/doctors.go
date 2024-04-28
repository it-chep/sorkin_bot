package appointment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sorkin_bot/internal/domain/services/appointment/mis_dto"
)

func (as *AppointmentService) GetDoctors(ctx context.Context) (err error) {
	op := "sorkin_bot.internal.domain.services.appointment.doctors.GetDoctors"
	var Request mis_dto.GetUserRequest
	var Response mis_dto.GetUsersResponse

	Request = mis_dto.GetUserRequest{}
	jsonBody, err := json.Marshal(Request)
	if err != nil {
		as.logger.Error(fmt.Sprintf("error while marshalling json %s \nplace: %s", err, op))
		return err
	}
	body := bytes.NewReader(jsonBody)
	responseBody := as.sendToMIS(ctx, mis_dto.GetUsersMethod, body)

	err = json.Unmarshal(responseBody, &Response)
	if err != nil {
		as.logger.Info(fmt.Sprintf("error while unmarshalling json %s \nplace: %s", err, op))
		return
	}

	// Вывод данных о докторах
	for _, doctor := range Response.Data {
		fmt.Printf("ID: %d\n", doctor.ID)
		fmt.Printf("Имя: %s\n", doctor.Name)
		fmt.Printf("Дата рождения: %s\n", doctor.BirthDate)
		fmt.Printf("Роль: %s\n", doctor.RoleTitles)
		fmt.Println("--------------------------------------")
	}

	return nil
}

func (as *AppointmentService) ChooseDoctors(ctx context.Context) {

}

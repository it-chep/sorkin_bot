package adapter

import (
	"context"
	"fmt"
	"sorkin_bot/internal/clients/gateways/dto"
	entity "sorkin_bot/internal/domain/entity/user"
)

func (a *AppointmentServiceAdapter) GetPatientById(ctx context.Context, patientId int) (patientDTO dto.CreatedPatientDTO, err error) {
	patientDTO, err = a.gateway.GetPatientById(ctx, patientId)

	if err != nil {
		return dto.CreatedPatientDTO{}, err
	}
	return patientDTO, nil
}

func (a *AppointmentServiceAdapter) CreatePatient(ctx context.Context, user entity.User) (patientId *int, err error) {
	patientDTO, err := a.gateway.GetPatientByBirthDate(ctx, user)
	fmt.Println("CreatePatient", patientDTO, &patientDTO.PatientID)
	if err == nil {
		return &patientDTO.PatientID, nil
	}

	userDTO := dto.PatientDTO{
		LastName:  *user.GetLastName(),
		FirstName: user.GetFirstName(),
		ThirdName: user.GetThirdName(),
		BirthDate: *user.GetBirthDate(),
		Phone:     *user.GetPhone(),
	}
	patientId, err = a.gateway.CreatePatient(ctx, userDTO)

	if err != nil {
		return nil, err
	}
	return patientId, nil
}

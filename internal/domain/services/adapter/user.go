package adapter

import (
	"context"
	"sorkin_bot/internal/clients/gateways/dto"
)

func (a *AppointmentServiceAdapter) GetPatientById(ctx context.Context, patientId int) (patientDTO dto.CreatedPatientDTO, err error) {
	patientDTO, err = a.gateway.GetPatientById(ctx, patientId)

	if err != nil {
		return dto.CreatedPatientDTO{}, err
	}
	return patientDTO, nil
}

func (a *AppointmentServiceAdapter) CreatePatient(ctx context.Context, userDTO dto.PatientDTO) (patientId *int, err error) {
	patientId, err = a.gateway.CreatePatient(ctx, userDTO)
	if err != nil {
		return nil, err
	}
	return patientId, nil
}

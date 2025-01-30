package notification

import (
	"context"
	"sorkin_bot/internal/clients/gateways/dto"
)

type notifyGateway interface {
	SendNotification(ctx context.Context, to []string, message string) error
}

type misGateway interface {
	GetPatientById(ctx context.Context, patientId int) (patientDTO dto.CreatedPatientDTO, err error)
}

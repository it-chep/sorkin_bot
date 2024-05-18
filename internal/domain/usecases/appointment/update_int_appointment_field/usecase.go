package update_int_appointment_field

import (
	"context"
	"log/slog"
)

type UpdateIntAppointmentFieldUseCase struct {
	writeRepo writeRepo
	logger    *slog.Logger
}

func NewUpdateIntAppointmentFieldUseCase(writeRepo writeRepo, logger *slog.Logger) UpdateIntAppointmentFieldUseCase {
	return UpdateIntAppointmentFieldUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc UpdateIntAppointmentFieldUseCase) Execute(ctx context.Context, tgId int64, fieldValue int, fieldName string) error {
	return uc.writeRepo.UpdateIntDraftAppointment(ctx, tgId, fieldValue, fieldName)
}

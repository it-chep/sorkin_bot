package update_str_appointment_field

import (
	"context"
	"log/slog"
)

type UpdateStrAppointmentFieldUseCase struct {
	writeRepo writeRepo
	logger    *slog.Logger
}

func NewUpdateStrAppointmentFieldUseCase(writeRepo writeRepo, logger *slog.Logger) UpdateStrAppointmentFieldUseCase {
	return UpdateStrAppointmentFieldUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc UpdateStrAppointmentFieldUseCase) Execute(ctx context.Context, tgId int64, fieldValue, fieldName string) error {
	return uc.writeRepo.UpdateStrFieldDraftAppointment(ctx, tgId, fieldValue, fieldName)
}

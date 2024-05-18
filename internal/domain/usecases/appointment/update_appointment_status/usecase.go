package update_appointment_status

import (
	"context"
	"log/slog"
)

type UpdateAppointmentStatusUseCase struct {
	writeRepo writeRepo
	logger    *slog.Logger
}

func NewUpdateAppointmentStatusUseCase(writeRepo writeRepo, logger *slog.Logger) UpdateAppointmentStatusUseCase {
	return UpdateAppointmentStatusUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc UpdateAppointmentStatusUseCase) Execute(ctx context.Context, tgId int64) error {
	return uc.writeRepo.UpdateStatusDraftAppointment(ctx, tgId)
}

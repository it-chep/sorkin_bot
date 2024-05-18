package clean_draft_appointment

import (
	"context"
	"log/slog"
)

type CleanDraftAppointmentUseCase struct {
	writeRepo writeRepo
	logger    *slog.Logger
}

func NewCleanDraftAppointmentUseCase(writeRepo writeRepo, logger *slog.Logger) CleanDraftAppointmentUseCase {
	return CleanDraftAppointmentUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc CleanDraftAppointmentUseCase) Execute(ctx context.Context, tgId int64) error {
	return uc.writeRepo.CleanDraftAppointment(ctx, tgId)
}

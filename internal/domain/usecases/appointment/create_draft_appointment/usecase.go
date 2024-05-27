package create_draft_appointment

import (
	"context"
	"log/slog"
)

type CreateDraftAppointmentUseCase struct {
	writeRepo writeRepo
	logger    *slog.Logger
}

func NewCreateDraftAppointmentUseCase(writeRepo writeRepo, logger *slog.Logger) CreateDraftAppointmentUseCase {
	return CreateDraftAppointmentUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc CreateDraftAppointmentUseCase) Execute(ctx context.Context, tgId int64) error {
	return uc.writeRepo.CreateEmptyDraftAppointment(ctx, tgId)
}

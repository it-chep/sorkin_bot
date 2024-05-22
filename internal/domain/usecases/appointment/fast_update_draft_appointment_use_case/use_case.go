package fast_update_draft_appointment_use_case

import (
	"context"
	"log/slog"
	"sorkin_bot/internal/domain/entity/appointment"
)

type FastUpdateDraftAppointmentUseCase struct {
	writeRepo writeRepo
	logger    *slog.Logger
}

func NewFastUpdateDraftAppointmentUseCase(writeRepo writeRepo, logger *slog.Logger) FastUpdateDraftAppointmentUseCase {
	return FastUpdateDraftAppointmentUseCase{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc FastUpdateDraftAppointmentUseCase) Execute(ctx context.Context, tgId int64, draftAppointment appointment.DraftAppointment) error {
	return uc.writeRepo.FastUpdateDraftAppointment(ctx, tgId, draftAppointment)
}

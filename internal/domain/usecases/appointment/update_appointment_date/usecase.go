package update_appointment_date

import (
	"context"
	"log/slog"
)

type UpdateAppointmentDate struct {
	writeRepo writeRepo
	logger    *slog.Logger
}

func NewUpdateAppointmentDate(writeRepo writeRepo, logger *slog.Logger) UpdateAppointmentDate {
	return UpdateAppointmentDate{
		writeRepo: writeRepo,
		logger:    logger,
	}
}

func (uc UpdateAppointmentDate) Execute(ctx context.Context, tgId int64, timeStart, timeEnd, date string) error {
	return uc.writeRepo.UpdateDateDraftAppointment(ctx, tgId, timeStart, timeEnd, date)
}

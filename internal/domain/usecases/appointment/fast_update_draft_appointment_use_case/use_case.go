package fast_update_draft_appointment_use_case

import (
	"context"
	"fmt"
	"log/slog"
	"sorkin_bot/internal/domain/entity/appointment"
	"strings"
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

func (uc FastUpdateDraftAppointmentUseCase) Execute(ctx context.Context, tgId int64, draftAppointment appointment.DraftAppointment, created bool) error {
	if created {
		//todo по хорошему вынести в 1 update чтобы снизить нагрузку на базу
		err := uc.writeRepo.UpdateIntDraftAppointment(ctx, tgId, *draftAppointment.GetDoctorId(), "doctor_id")
		if err != nil {
			return err
		}

		if draftAppointment.GetTimeStart() != nil {
			err = uc.writeRepo.UpdateDateDraftAppointment(
				ctx, tgId, *draftAppointment.GetTimeStart(),
				*draftAppointment.GetTimeEnd(),
				strings.Split(*draftAppointment.GetTimeStart(), " ")[0],
			)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("timeStart is nil")
		}
		return nil
	} else {
		return uc.writeRepo.FastUpdateDraftAppointment(ctx, tgId, draftAppointment)
	}
}

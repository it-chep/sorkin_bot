package fast_update_draft_appointment_use_case

import (
	"context"
	"sorkin_bot/internal/domain/entity/appointment"
)

type writeRepo interface {
	FastUpdateDraftAppointment(
		ctx context.Context, tgId int64,
		draftAppointment appointment.DraftAppointment,
	) (err error)

	UpdateDateDraftAppointment(
		ctx context.Context, tgId int64, timeStart, timeEnd, date string,
	) (err error)
	UpdateIntDraftAppointment(ctx context.Context, tgId int64, intValue int, intField string) (err error)
}

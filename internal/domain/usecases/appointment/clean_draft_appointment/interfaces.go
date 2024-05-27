package clean_draft_appointment

import "context"

type writeRepo interface {
	CleanDraftAppointment(ctx context.Context, tgId int64) (err error)
}

package create_draft_appointment

import "context"

type writeRepo interface {
	CreateEmptyDraftAppointment(ctx context.Context, tgId int64) (err error)
}

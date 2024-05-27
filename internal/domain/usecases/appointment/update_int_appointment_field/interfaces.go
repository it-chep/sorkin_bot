package update_int_appointment_field

import "context"

type writeRepo interface {
	UpdateIntDraftAppointment(ctx context.Context, tgId int64, intValue int, intField string) (err error)
}

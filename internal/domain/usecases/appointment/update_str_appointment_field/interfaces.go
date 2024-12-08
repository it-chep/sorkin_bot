package update_str_appointment_field

import "context"

type writeRepo interface {
	UpdateStrFieldDraftAppointment(ctx context.Context, tgId int64, strValue, strField string) (err error)
}

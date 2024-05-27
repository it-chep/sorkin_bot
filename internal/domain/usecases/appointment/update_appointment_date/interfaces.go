package update_appointment_date

import "context"

type writeRepo interface {
	UpdateDateDraftAppointment(
		ctx context.Context, tgId int64, timeStart, timeEnd, date string,
	) (err error)
}
